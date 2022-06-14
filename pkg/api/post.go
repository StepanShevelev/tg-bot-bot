package api

import (
	mydb "github.com/StepanShevelev/tg-bot-bot/pkg/db"
	"github.com/sirupsen/logrus"
	"html/template"
	"net/http"
	"strconv"
)

func apiGetPost(w http.ResponseWriter, r *http.Request) {
	if !isMethodGET(w, r) {
		return
	}

	postName, okId := ParseName(w, r)
	if !okId {
		w.Write([]byte(`{"error": "can't pars name"}`))
		return
	}
	//fmt.Println(postName)

	GetHtmlByName(postName, w)
	//if err != nil {
	//	mydb.UppendErrorWithPath(err)
	//	logrus.Info("Can`t find post", err)
	//	w.WriteHeader(http.StatusNotFound)
	//
	//}
	//SendData(&result, w)
}

func GetPostByName(postName string) (*mydb.Post, error) {
	var post *mydb.Post

	result := mydb.Database.Db.Preload("Images").Find(&post, "Name = ?", postName)
	if result.Error != nil {
		//w.WriteHeader(http.StatusNotFound)
		mydb.UppendErrorWithPath(result.Error)
		logrus.Info("Can`t find post", result.Error)
		return nil, result.Error
	}

	return post, nil

}

func GetImagesByPost(postId uint) ([]mydb.Image, error) {

	var images []mydb.Image

	result := mydb.Database.Db.Find(&images, "post_id = ?", postId)
	if result.Error != nil {
		mydb.UppendErrorWithPath(result.Error)
		logrus.Info("Could not find post", result.Error)
		return nil, result.Error
	}
	return images, nil

}

func GetHtmlByName(postName string, w http.ResponseWriter) {

	post, err := GetPostByName(postName)
	if err != nil {
		mydb.UppendErrorWithPath(err)
		logrus.Info("Could not find post to create HTML", err)
		return
	}

	images, err := GetImagesByPost(post.ID)
	if err != nil {
		mydb.UppendErrorWithPath(err)
		logrus.Info("Could not find post to create HTML", err)
		return
	}

	var tmpl = `
	<html>
<head>
<title>

<h1 class={{$.Title}}></h1>

</title>
</head>
<body>



<div>{{$.Title}}</div>


<div><img src={{  printf "%s" ((index  .Images 0).Name )  }} style="width: 500px;height: 400px" ></div>


{{.Text}}

{{range .Images}}
{{if index 3}}{{break}}{{end}}
{{continue}}

<img src={{ .Name}}>
{{end}}
</body>
</html>
`
	a := mydb.Post{
		Title:  post.Title,
		Text:   post.Text,
		Images: images,
	}

	for i, _ := range images {
		if i == 0 {
			continue
		}
		tmpl = tmpl + "<img src= {{ btoa (index $.Images" + " " + strconv.Itoa(i) + ").Name }} style=\"width: 500px;height: 400px\" >"
	}

	// Make and parse the HTML template
	t, err := template.New(post.Title).Funcs(template.FuncMap{
		"btoa": func(b []byte) string { return string(b) },
	}).Parse(tmpl)
	if err != nil {
		logrus.Info("Error occurred while creating new template", err)
		mydb.UppendErrorWithPath(err)
		return
	}
	err = t.Execute(w, a)
	if err != nil {
		logrus.Info("Error occurred while updating file data", err)
		mydb.UppendErrorWithPath(err)
		return
	}
	//file, err := os.Create(post.Title + ".html")
	//if err != nil {
	//	logrus.Info("Error occurred while creating file", err)
	//	mydb.UppendErrorWithPath(err)
	//	return "", err
	//}

	//err = t.ExecuteTemplate(file, post.Title, a)
	//if err != nil {
	//	logrus.Info("Error occurred while updating file data", err)
	//	mydb.UppendErrorWithPath(err)
	//	return "", err
	//}
	//
	////post.WhoTookMe = whoTookMe
	//result := mydb.Database.Db.Save(&post)
	//if result.Error != nil {
	//	logrus.Info("Error occurred while updating post", err)
	//	mydb.UppendErrorWithPath(result.Error)
	//}
	//replacer := strings.NewReplacer(" ", "", "/", "", ".", "", ",", "", "!", "", ":", "", "?", "")
	//txt := replacer.Replace(post.Title)
	//path := "./" + txt + ".html"
	//return path, nil
}
