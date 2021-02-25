package main

import (
	"flag"
	pb "github.com/Ghost-Pacer/input-goi2c/experiments/bno055live/proto"
	"github.com/Ghost-Pacer/input-goi2c/pkg/bno055"
	"github.com/go-zeromq/zmq4"
	"github.com/golang/protobuf/ptypes"
	"google.golang.org/protobuf/proto"
	"gopkg.in/tomb.v2"
	"log"
	"os"
	"os/signal"
	"periph.io/x/conn/v3/i2c/i2creg"
	"periph.io/x/host/v3"
	"syscall"
	"time"
)

// protoc invocation:
// C:\Users\Jensen Turner\GitProjects\input-goi2c\experiments\bno055live>protoc --go_out=. --go_opt=module=github.com/Ghost-Pacer/input-goi2c/experiments/bno055live --proto_path=..\..\..\protocols\ghostpacer\input i2c_devices.proto

var I2CBus = flag.String("b", "2", "I2C bus")
var I2CAddr = flag.Int("a", 0x28, "I2C address")
var RefreshInterval = flag.Duration("r", 8*time.Millisecond, "refresh interval")
var SocketEndpoint = flag.String("e", "tcp://localhost:51101", "ZMQ bound endpoint")
var Verbose = flag.Bool("v", false, "print all sent messages (do not use in production)")

func printv(vals ...interface{}) {
	if *Verbose {
		log.Println(vals...)
	}
}

func run() error {
	mainTomb := new(tomb.Tomb)

	socket := zmq4.NewPub(mainTomb.Context(nil))
	defer socket.Close()

	if err := socket.Listen(*SocketEndpoint); err != nil {
		return err
	}
	/*if err := socket.SetOption("CONFLATE", true); err != nil {
		return err
	}*/
	log.Println("zmq4: listening on", *SocketEndpoint)

	if _, err := host.Init(); err != nil {
		return err
	}
	log.Println("periph: initted host")

	bus, err := i2creg.Open(*I2CBus)
	if err != nil {
		return err
	}
	log.Println("periph: initted bus")
	defer bus.Close()

	bno, err := bno055.New(bus, uint16(*I2CAddr))
	log.Println("bno055: initted bno055")
	if err != nil {
		return err
	}
	defer bno.Halt()

	signalChannel := make(chan os.Signal, 1)
	signal.Notify(signalChannel, syscall.SIGINT, syscall.SIGTERM)

	mainTomb.Go(func() error {
		return doWork(mainTomb, socket, bno)
	})

	caughtSignal := <-signalChannel
	log.Println("Caught", caughtSignal, "shutting down...")
	mainTomb.Kill(nil)
	return mainTomb.Wait()
}

func doWork(tomb *tomb.Tomb, socket zmq4.Socket, bno *bno055.Dev) error {
	ticker := time.NewTicker(*RefreshInterval)
	defer ticker.Stop()

	for {
		select {
		case <-tomb.Dying():
			return nil
		case <-ticker.C:
			printv("ticked")

			sourced := ptypes.TimestampNow()

			quat, err := bno.ReadQuat()
			if err != nil {
				return err
			}
			printv("\tgot quat", quat)

			eul, err := bno.ReadEuler()
			if err != nil {
				return err
			}
			printv("\tgot eul", eul)

			lin, err := bno.ReadLinearAccel()
			if err != nil {
				return err
			}
			printv("\tgot lin", lin)

			snapshot := &pb.IMUSnapshot{
				EventTimings: &pb.InputEventTimings{
					Sourced: sourced,
					Updated: ptypes.TimestampNow(),
				},
				OrientationQuaternion: &pb.Vec4{
					W: quat[0],
					X: quat[1],
					Y: quat[2],
					Z: quat[3],
				},
				LinearAcceleration: &pb.Vec3{
					X: lin[0],
					Y: lin[1],
					Z: lin[2],
				},
			}
			out, err := proto.Marshal(snapshot)
			if err != nil {
				return err
			}

			if err := socket.Send(zmq4.NewMsg(out)); err != nil {
				return err
			}
			printv("\tsent on socket")
		}
	}
}

func main() {
	flag.Parse()
	log.SetFlags(log.LstdFlags | log.Lmicroseconds)
	if err := run(); err != nil {
		log.Fatal(err)
	}
}
