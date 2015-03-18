package jobs

import (
  "fmt"
  "io/ioutil"
  "net/http"
  "unsafe"
  "strings"
  "strconv"
  "github.com/moovweb/gokogiri/xml"
  "github.com/moovweb/gokogiri/xpath"
)

import (
  "github.com/revel/revel"
  "github.com/revel/revel/modules/jobs/app/jobs"
  "time"
)

import "appratings/app/models"

type UpdateRatingsJob struct {
}

func (j UpdateRatingsJob) DownloadURL(url string, user_agent string) ([]byte, error) {
  client := &http.Client {}

  req, err := http.NewRequest("GET", url, nil)
  req.Header.Add("User-Agent", user_agent)

  resp, err := client.Do(req)
  if err != nil {
    return nil, err
  }
  
  bytes, err := ioutil.ReadAll(resp.Body)
  if err != nil {
    return nil, err
  }

  resp.Body.Close()

  return bytes, nil
}

type StarRatingMap map[int64]int64

func (j UpdateRatingsJob) ParseStarRatings(data []byte) (total, version StarRatingMap, err error) {
  xmlDoc, err := xml.Parse(data, xml.DefaultEncodingBytes, nil, xml.DefaultParseOption, xml.DefaultEncodingBytes)

  if err != nil {
    return nil, nil, err
  } else {
    xmlDoc.RecursivelyRemoveNamespaces()

    xp := xmlDoc.DocXPathCtx()
    expr := xpath.Compile("//HBoxView[@rightInset=5]/@alt")

    nodePtr := unsafe.Pointer(xmlDoc.Root())
    err := xp.Evaluate(nodePtr, expr)

    if err != nil {
      return nil, nil, err
    } else {
      retType := xp.ReturnType()

      if (retType == 1) {
        nodes, err := xp.ResultAsNodeset()

        if err != nil {
          return nil, nil, err
        } else {
          total := make(StarRatingMap)
          version := make(StarRatingMap)

          for index, element := range nodes {
            node := xml.NewNode(element, xmlDoc)
            nodeStr := node.String()
            nodeSplitString := strings.Split(nodeStr, ", ")

            stars := nodeSplitString[0]
            ratings := nodeSplitString[1]

            star, _ := strconv.ParseInt(strings.Split(stars, " ")[0], 10, 0)
            rating, _ := strconv.ParseInt(strings.Split(ratings, " ")[0], 10, 0)

            if (index < 5) {
              total[star] = rating
            } else {
              version[star] = rating
            }
          }

          return total, version, nil
        }
      }
    }

    return nil, nil, nil
  }
}

func (j UpdateRatingsJob) DownloadAppRatings(appId string) (total, version StarRatingMap, err error) {
  url := fmt.Sprintf("https://itunes.apple.com/WebObjects/MZStore.woa/wa/viewContentsUserReviews?id=%s&pageNumber=0&sortOrdering=2&type=Purple+Software&ign-mscache=1", appId)

  userAgent := "iTunes/11.0.2 (Macintosh; OS X 10.8.2) AppleWebKit/536.26.14"
  xmlBytes, err := j.DownloadURL(url, userAgent)

  if err != nil {
    return nil, nil, err
  } else {
    return j.ParseStarRatings(xmlBytes)
  }
}

func (j UpdateRatingsJob) Update() {
  cn := make(chan int)

  apps := models.FindAllApps()

  for index, app := range apps {
    go func(appId string) {
      _, _, err := j.DownloadAppRatings(appId)
    
      if err != nil {
        revel.ERROR.Printf("Error fetching app ratings: %s", err)
      } else {
        revel.INFO.Printf("Success fetching ratings for app: %s", appId)
      }

      cn <- index
    }(app.AppID)
  }

  for i := 0; i < len(apps); i++ {
    <- cn
  }
}

func (j UpdateRatingsJob) Run() {
  j.Update()
}

func init() {
  revel.OnAppStart(func() {
    jobs.Now(UpdateRatingsJob{})
    jobs.Every(1 * time.Hour, UpdateRatingsJob{})
  })
}
