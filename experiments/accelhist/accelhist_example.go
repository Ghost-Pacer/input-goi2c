package accelhist

import (
	"github.com/Ghost-Pacer/input-goi2c/pkg/nearsight"
	"gonum.org/v1/gonum/mat"
	"gonum.org/v1/gonum/num/quat"
	"gonum.org/v1/gonum/spatial/r3"
)

const AccelHistLength = 50

var _ nearsight.MatVecDenseStack = nearsight.NewMatVecDenseStack(AccelHistLength)

func UpdateAccelHist(AccelHist nearsight.MatVecDenseStack, quat quat.Number, linaccel r3.Vec) {
	newAccelEntry := mat.NewVecDense(3, []float64{
		quat.Real*linaccel.Y - quat.Kmag*linaccel.Z + quat.Jmag*linaccel.X,
		quat.Real*linaccel.Z - quat.Imag*linaccel.X + quat.Kmag*linaccel.Y,
		quat.Real*linaccel.X - quat.Jmag*linaccel.Y + quat.Imag*linaccel.Z,
	})
	AccelHist.Push(*newAccelEntry)
}

func UseAccelHist(AccelHist nearsight.MatVecDenseStack) {
	_ = AccelHist
}
