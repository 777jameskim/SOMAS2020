package team5

import (
	"math"

	"github.com/SOMAS2020/SOMAS2020/internal/common/shared"
)

type forecastInfo struct {
	epiX       shared.Coordinate // x co-ord of disaster epicentre
	epiY       shared.Coordinate // y ""
	mag        shared.Magnitude
	turn       uint
	confidence float64
}

type forecastHistory map[uint]forecastInfo // stores history of past disasters

// MakeDisasterPrediction is called on each client for them to make a prediction about a disaster
// Prediction includes location, magnitude, confidence etc
// COMPULSORY, you need to implement this method
func (c *client) MakeDisasterPrediction() shared.DisasterPredictionInfo {

	meanDisaster := c.getMeanDisasterInfo()
	prediction := shared.DisasterPrediction{
		CoordinateX: meanDisaster.epiX,
		CoordinateY: meanDisaster.epiY,
		Magnitude:   meanDisaster.mag,
		TimeLeft:    int(meanDisaster.turn - c.getTurn()),
	}

	prediction.Confidence = c.determineForecastConfidence()
	trustedIslandIDs := []shared.ClientID{}
	trustThresh := c.config.forecastTrustTreshold
	for id := range c.getTrustedTeams(trustThresh, false, forecastingBasis) {
		trustedIslandIDs = append(trustedIslandIDs, id)
	}

	// Return all prediction info and store our own island's prediction in global variable
	predictionInfo := shared.DisasterPredictionInfo{
		PredictionMade: prediction,
		TeamsOfferedTo: trustedIslandIDs,
	}
	c.lastDisasterPrediction = prediction
	// update forecast history
	c.forecastHistory[c.getTurn()] = forecastInfo{
		epiX: prediction.CoordinateX,
		epiY: prediction.CoordinateY,
		mag:  prediction.Magnitude,
		turn: uint(prediction.TimeLeft) + c.getTurn(),
	}
	return predictionInfo
}

// averages observations over history to get 'mean' disaster
func (c client) getMeanDisasterInfo() forecastInfo {
	sumX, sumY, sumMag := 0.0, 0.0, 0.0

	for _, dInfo := range c.disasterHistory {
		sumX += dInfo.report.X
		sumY += dInfo.report.Y
		sumMag += dInfo.report.Y
	}
	n := float64(len(c.forecastHistory))
	period, conf := c.analyseDisasterPeriod()

	meanDisaster := forecastInfo{
		epiX:       sumX / n,
		epiY:       sumY / n,
		mag:        sumMag / n,
		turn:       c.getLastDisasterTurn() + period,
		confidence: conf,
	}
	return meanDisaster
}

func (c client) getLastDisasterTurn() uint {
	n = len(c.disasterHistory)
	if n > 0 {
		lastT := uint(0)
		for t := range c.disasterHistory { // TODO: find nicer way of getting largest key (turn)
			if t > lastT {
				lastT = t
			}
		}
		return lastT
	}
	return 0
}

func (c *client) analyseDisasterPeriod() (period uint, conf float64) {
	periods := []uint{0}
	periodDiffs := []int{}
	i := 1
	for turn := range c.disasterHistory {
		periods = append(periods, turn-periods[i-1]) // period = no. turns between successive disasters
		if len(periods) > 2 {
			periodDiffs = append(periodDiffs, int(periods[i]-periods[i-1]))
		}
		i++
	}
	periods = periods[1:] // remove leading 0
	diffSum := 0.0        // sum of differences in periods between disasters
	for _, pd := range periodDiffs {
		diffSum += math.Abs(float64(pd))
	}
	if diffSum == 0 { // perfectly cyclical - consistent period
		return periods[0], 100.0
	}
	// if not consisten, return mean period we've seen so far
	return uint(diffSum / float64(len(periodDiffs))), 65.0
}

func (c *client) determineForecastConfidence() float64 {
	totalDisaster := forecastInfo{}
	sqDiff := func(x, meanX float64) float64 { return math.Pow(x-meanX, 2) }
	meanInfo := c.getMeanDisasterInfo()
	// Find the sum of the square of the difference between the actual and mean, for each field
	for _, d := range c.forecastHistory {
		totalDisaster.epiX += sqDiff(d.epiX, meanInfo.epiX)
		totalDisaster.epiY += sqDiff(d.epiY, meanInfo.epiY)
		totalDisaster.mag += sqDiff(d.mag, meanInfo.mag)
	}

	// TODO: find a better method of calculating confidence
	// Find the sum of the variances and the average variance
	variance := (totalDisaster.epiX + totalDisaster.epiY + totalDisaster.mag) / float64(len(c.forecastHistory))
	variance = math.Min(c.config.maxForecastVariance, variance)

	return c.config.maxForecastVariance - variance
}

// ReceiveDisasterPredictions provides each client with the prediction info, in addition to the source island,
// that they have been granted access to see
// COMPULSORY, you need to implement this method
func (c *client) ReceiveDisasterPredictions(receivedPredictions shared.ReceivedDisasterPredictionsDict) {
	// If we assume that we trust each island equally (including ourselves), then take the final prediction
	// of disaster as being the weighted mean of predictions according to confidence

	sumX, sumY, sumMag, sumConf := 0.0, 0.0, 0.0, 0.0
	sumTime := 0

	//c.lastDisasterForecast.Confidence *= 1.3 // inflate confidence of our prediction above others
	receivedPredictions[ourClientID] = shared.ReceivedDisasterPredictionInfo{PredictionMade: c.lastDisasterPrediction, SharedFrom: ourClientID}

	// weight predictions by their confidence and our assessment of their forecasting reputation
	for rxTeam, pred := range receivedPredictions {
		rep := float64(c.opinions[rxTeam].forecastReputation) + 1 // our notion of another island's forecasting reputation
		sumX += pred.PredictionMade.Confidence * pred.PredictionMade.CoordinateX * rep
		sumY += pred.PredictionMade.Confidence * pred.PredictionMade.CoordinateY * rep
		sumMag += pred.PredictionMade.Confidence * pred.PredictionMade.Magnitude * rep
		sumTime += int(pred.PredictionMade.Confidence) * pred.PredictionMade.TimeLeft * int(rep)
		sumConf += pred.PredictionMade.Confidence * rep
	}

	// Finally get the final prediction generated by considering predictions from all islands that we have available
	finalPrediction := shared.DisasterPrediction{
		CoordinateX: sumX / sumConf,
		CoordinateY: sumY / sumConf,
		Magnitude:   sumMag / sumConf,
		TimeLeft:    int((float64(sumTime) / sumConf) + 0.5), // +0.5 for rounding
		Confidence:  sumConf / float64(len(receivedPredictions)),
	}

	c.Logf("Final Prediction: [%v]", finalPrediction)
}

func (c *client) analysePastDisasters() {

}
