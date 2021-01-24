package nearsight

import (
	"github.com/Ghost-Pacer/input-goi2c/nearsight/cow"
	"gonum.org/v1/gonum/mat"
	"gonum.org/v1/gonum/num/quat"
	"gonum.org/v1/gonum/spatial/r3"
)

const AccelHistLength = 50

var AccelHist *cow.MatVecDenseCircularQueue = cow.NewMatVecDenseCircularQueue(AccelHistLength)

func UpdateAccelHist(AccelHist *cow.MatVecDenseCircularQueue, quat quat.Number, linaccel r3.Vec) {
	newAccelEntry := mat.NewVecDense(3, []float64{
		quat.Real*linaccel.Y - quat.Kmag*linaccel.Z + quat.Jmag*linaccel.X,
		quat.Real*linaccel.Z - quat.Imag*linaccel.X + quat.Kmag*linaccel.Y,
		quat.Real*linaccel.X - quat.Jmag*linaccel.Y + quat.Imag*linaccel.Z,
	})
	AccelHist.Push(*newAccelEntry)
}

func UseAccelHist(AccelHist *cow.MatVecDenseCircularQueue) {
	// AccelHist is just a pointer to a slice but also has convenience methods defined on MatVecDenseRing
}
