package facebook

import (
  "strings"
  "log"
  "io/ioutil"
  "errors"
  "encoding/json"
  "net/http"
  "strconv"
)

type App struct {
  AppID string
  AppSecret string
}

type User struct {
    Name string `json:"name"`
    Username string `json:"username,omitEmpty"`
    Currency string `json:"currency"`
    FacebookId int64 `json:"id"`
    User_Id string `json:"user_id"`
    Gender string `json:"gender"`
    Age int `json:"age"`
    Birthday string `json:"birthday,omitEmpty"` 
    Location string `json:"location"`
    Attributes []byte `json:"attributes,omitEmpty"`
}


type FacebookAuthResponse struct {
    AccessToken string `schema:"accessToken"`
    ExpiresIn int `schema:"expiresIn"`
    SignedRequest string `schema:"signedRequest"`
    UserId string `schema:"userId"`
}

type FacebookTokenDebugInformation struct {
    Data TokenDebugData `json:"data"`
}

type TokenDebugData struct {
    AppId string `json:"app_id"`
    UserId string  `json:"user_id"`

}

type FacebookUser struct {
    Name string `json:"name"`
    Currency interface{} `json:"currency,omitEmpty"`
    Id string `json:"id"`
    Gender string `json:"gender"`
    Birthday string `json:"birthday,omitEmpty"` 
    Location interface{} `json:"location"`
}


func (fbu *FacebookUser) ToGiftedUser() (User,error){
    var u User

    u.Name = fbu.Name
    u.Currency = fbu.Currency.(map[string]interface{})["user_currency"].(string)
    u.Gender = fbu.Gender
    u.Location = fbu.Location.(map[string]interface{})["name"].(string)

    err := fbu.ParseBirthdayString()
    if err != nil {
        return u, err
    }

    u.Birthday = fbu.Birthday

    convertedId, errConv := strconv.ParseInt(fbu.Id, 10, 64)
    u.FacebookId = convertedId
    if errConv != nil {
        return u, errConv
    }

    return u, nil
}

func (fbu *FacebookUser) ParseBirthdayString() error {
    //parse facebook date into "YYYY-MM-DD" formatting, also deal users who do not share all birthday detailsa
    split := strings.Split(fbu.Birthday ,"/")
    var err error
    switch {
    case len(split)==3:
        fbu.Birthday = split[2]+"-"+split[0]+"-"+split[1]
    case len(split)==2:
        fbu.Birthday = "1000"+"-"+split[0]+"-"+split[1]
    case len(split)==1:
        fbu.Birthday = fbu.Birthday
    default:
        err = errors.New("Error: Facebook provided invalid date format")
    }
    return err
}

func (app *App) CheckAccessToken(authResponse *FacebookAuthResponse) (bool, error){
    var url string
    url = fmt.Sprintf("https://graph.facebook.com/debug_token?input_token=%s&access_token=%s|%s" , authResponse.AccessToken, facebookAppID, facebookClientSecret)

    resp, err := http.Get(url)
    if err != nil {
        return false, err
    }

    body, bodyErr := ioutil.ReadAll(resp.Body)
    if bodyErr != nil {
        return false, bodyErr
    }

    debug := new(FacebookTokenDebugInformation)

    jsonErr := json.Unmarshal(body, &debug)
    if jsonErr != nil {
        return false, jsonErr
    }

    if debug.Data.UserId == authResponse.UserId && debug.Data.AppId == facebookAppID {
        return true, nil
    } 
  return false, errors.New("token inspection failed")
}

  
  
func (authResponse *FacebookAuthResponse) CheckAccessToken() (bool,error) {
    var url string
    url = fmt.Sprintf("https://graph.facebook.com/debug_token?input_token=%s&access_token=%s|%s" , authResponse.AccessToken, facebookAppID, facebookClientSecret)

    resp, err := http.Get(url)
    if err != nil {
        return false, err
    }

    body, bodyErr := ioutil.ReadAll(resp.Body)
    if bodyErr != nil {
        return false, bodyErr
    }

    debug := new(FacebookTokenDebugInformation)

    jsonErr := json.Unmarshal(body, &debug)
    if jsonErr != nil {
        return false, jsonErr
    }

    if debug.Data.UserId == authResponse.UserId && debug.Data.AppId == facebookAppID {
        return true, nil
    } 
  return false, errors.New("token inspection failed")
}
