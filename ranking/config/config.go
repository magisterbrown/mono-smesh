package config
//import "github.com/mafredri/go-trueskill"

var GameTag string = "game:latest"
var GameFolder string = "games/"

const (
    PDraw = 33.0
    DefMu = 700.0
    DefSig = DefMu / 3
)


//func TsConfig() trueskill.Config {
//    prob, _ := trueskill.DrawProbability(33)
//    
//    return trueskill.New(
//        trueskill.Mu(DefMu),
//        trueskill.Sigma(DefSig),
//        trueskill.Beta(DefSig/2),
//        trueskill.Tau(DefSig/100),
//        prob)
//}


