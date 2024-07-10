package endpoint

import "os"

var SRV string
var GetInfoUrl string

func init() {

  SRV = os.Getenv("HOST") + ":" + os.Getenv("PORT")
  if SRV == "" {
    panic("Переменная окружения HOST или PORT не установлены")
  }
  GetInfoUrl = "http://" + SRV + "/info"
}
