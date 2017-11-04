package main
import (
"github.com/torniker/infermedica"
    "encoding/json"
    
    "strconv"
    "log"
    "net/http"
  "strings"
  "errors"
  "bytes"
    "github.com/satori/go.uuid"
    "fmt"
    "net/http/httptest"
    "os"
    "crypto/tls"
    cors "github.com/heppu/simple-cors"
)
// "github.com/gorilla/mux" 
//"
// "reflect"
//"io/ioutil"
//lkhhjhkjh
var (
    // WelcomeMessage A constant to hold the welcome message
    WelcomeMessage = "Welcome to Medico we want to gather some information ,How old are you?"

    // sessions = {
    //   "uuid1" = Session{...},
    //   ...
//app key
   
     app = infermedica.NewApp("a3e97cd7", "af637361b0854d49d87b6af528c086dc", "")
    // }
//ignore_groups =true
    sessions = map[string]Session{}

    processor = sampleProcessor
    
)

type (
    // Session Holds info about a session
    Session map[string]interface{}
   
    // JSON Holds a JSON object
    JSON map[string]interface{}
    
    // Processor Alias for Process func
    Processor func(session Session, message string) (string, error)
)

type D struct {
    Mentions     []infermedica.Evidence   `json:"mentions"`
   
}


type Group struct {
    Ignore_groups   bool   `json:"ignore_groups"`
   
}
type DiagnosisReq struct {
    Sex       infermedica.Sex        `json:"sex"`
    Age       int        `json:"age"`
    Evidences []infermedica.Evidence `json:"evidence"`
    Extras     Group             `json:"extras"`
}
 

func ProcessFunc(p Processor) {
    processor = p
}
func sampleProcessor(session Session, messagee string) (string, error) {
//app := infermedica.NewApp("a3e97cd7", "af637361b0854d49d87b6af528c086dc", "")
    message:=strings.ToLower(messagee)

    // Make sure a history key is defined in the session which points to a slice of strings
    _, ageFound := session["age"]
    if !ageFound {
        //
        ai, err :=strconv.ParseInt(message, 10, 64)
        if err != nil || ai<0 || ai>200  {
 return  fmt.Sprintf(`<font color="red"> this is not a valid age %s !</font>`, message),nil
}
        session["age"]=message

       return fmt.Sprintf("So, what is your gender male or female ?"), nil
    }

     _, sexFound := session["sex"]
    if !sexFound {
        
        if message!="male" && message!="female" {
 return fmt.Sprintf(`<font color="red">this is not a valid gender %s !</font>`, message),nil;
}
        
        session["sex"]=message

    return fmt.Sprintf("So, what is the certainity would you prefer between 0 and 1 ? "), nil

 

    }

_, probFound := session["probability"]
    if !probFound {
        i, err :=strconv.ParseFloat(message, 64)
        if err != nil || i > 1 || i < 0 {
 return fmt.Sprintf(`<font color="red">this is not a valid Probability %s !</font>`, message),nil;
}
 session["probability"]=message
return fmt.Sprintf("So, what do you suffer from?"), nil
}

_, evFound := session["Evidence"]
    if !evFound {
        err:=Parseei(message,session)
 if err != nil {
 return  fmt.Sprintf(`<font color="red"> this is not a valid symptom %s !</font>`, message),nil
}

var D DiagnosisReq


 x,_:=strconv.ParseInt(session["age"].(string), 10, 64)
// var x int64=session["age"].(int64)
if session["sex"].(string)=="male"{
  //reflect.ValueOf(session["age"]).Int()
  
 D=DiagnosisReq{Sex:infermedica.SexMale,Age:int(x) ,Evidences:session["Evidence"].([]infermedica.Evidence),Extras:Group{Ignore_groups:true}}
}else {

D=DiagnosisReq{Sex:infermedica.SexFemale,Age:int(x),Evidences:session["Evidence"].([]infermedica.Evidence),Extras:Group{Ignore_groups:true}}

}

m:=diagnosis(D)
session["questioninfo"]=m.Question.Items[0].ID


return  fmt.Sprintf(m.Question.Text),nil

    }

 
err2:=evadd(message,session)
if err2 != nil {
 return  fmt.Sprintf(`<font color="red"> this is not a valid choice  choose yes or no or unknown </font>`),nil
}
var D DiagnosisReq
var DD infermedica.DiagnosisReq



 x,_:=strconv.ParseInt(session["age"].(string), 10, 64)
// var x int64=session["age"].(int64)
if session["sex"].(string)=="male"{
  //reflect.ValueOf(session["age"]).Int()
   
 D=DiagnosisReq{Sex:infermedica.SexMale,Age:int(x) ,Evidences:session["Evidence"].([]infermedica.Evidence),Extras:Group{Ignore_groups:true}}
 DD=infermedica.DiagnosisReq{Sex:infermedica.SexMale,Age:int(x) ,Evidences:session["Evidence"].([]infermedica.Evidence)}
}else {

D=DiagnosisReq{Sex:infermedica.SexFemale,Age:int(x),Evidences:session["Evidence"].([]infermedica.Evidence),Extras:Group{Ignore_groups:true}}
DD=infermedica.DiagnosisReq{Sex:infermedica.SexMale,Age:int(x) ,Evidences:session["Evidence"].([]infermedica.Evidence)}

}

m:=diagnosis(D)

session["questioninfo"]=m.Question.Items[0].ID
maxx:=getmaximum(session,m.Conditions)
n,_:=strconv.ParseFloat(session["probability"].(string), 64)
if maxx>float64(n){
lab,_:=app.LabTestsRecommend(DD)
conddetaild:=getcondition(session)
 delete(session, "Evidence")

  var buffer bytes.Buffer
 var buffer2 bytes.Buffer
 if len(lab.Obligatory) == 0 {
buffer.WriteString(" No Obligatory Lab Test needed ")
 }else {
for i := 0; i < len(lab.Obligatory); i++ { 
 buffer.WriteString(lab.Obligatory[i].Name )
 buffer.WriteString("  ")

}}
if len(lab.Recommended) == 0{
buffer2.WriteString(" No Recommended Lab Test ")
}else{
for i := 0; i < len(lab.Recommended); i++ { 
 buffer2.WriteString(lab.Recommended[i].Name )
 buffer2.WriteString("  ")

}
}


return  fmt.Sprintf(`According to our analysis and your observation with probability: %f  you may have: %s known as: %s with  <br> </br>  
<strong>Prevalence</strong>: %s , <br> </br>
<strong>Acuteness</strong>: %s, <br> </br>
<strong>Severity</strong>: %s , <br> </br>
<strong>Hint</strong>: %s, <br> </br>
<strong>ICD10Code</strong>: %s, <br> </br>
<strong>TriageLevel</strong>: %s, <br> </br>
<strong>categories</strong>: %v, <br> </br>
<strong>Obligatory Labtests</strong>: %s, <br> </br>
<strong>Recommended Labtests</strong> %s, <br> </br>
<strong>if you are suffering from anything else Ask me :) , hope you get well soon </strong>`,
    maxx ,conddetaild.Condition.Name,conddetaild.Condition.CommonName,string(conddetaild.Prevalence),string(conddetaild.Acuteness),
    string(conddetaild.Severity),string(conddetaild.Extras.Hint),string(conddetaild.Extras.ICD10Code),string(conddetaild.TriageLevel),conddetaild.Categories,buffer.String(),
    buffer2.String()),nil
 }

return  fmt.Sprintf(m.Question.Text),nil
    }

func getcondition (session Session)(*infermedica.ConditionRes){

// app := infermedica.NewApp("a3e97cd7", "af637361b0854d49d87b6af528c086dc", "")
condition, _ := app.ConditionByID(session["maxcondid"].(string))
return condition
}
func getmaximum(session Session,cond []infermedica.DiagnosisConditionRes)(float64){

var max float64=-1;
var index int =0;
 for i := 0; i < len(cond); i++ { 


if cond [i].Probability>max {
    max=cond [i].Probability
    index=i
}

 }

session["maxcondid"]=cond [index].Condition.ID
return max

}

    // Add the message in the parsed body to the messages in the session
    //histo


    // Form a sentence out of the history in the form Message 1, Message 2, and Message 3
    //l := len(history)
    //wordsForSentence := make([]string, l)
    //copy(wordsForSentence, history)
    //if l > 1 {
    //    wordsForSentence[l-1] = "and " + wordsForSentence[l-1]
   // }
   // sentence := strings.Join(wordsForSentence, ", ")

    // Save the updated history to the session
    //session["history"] = history
func evadd (answer string,session Session)(error){
switch answer{
  case "yes": 
newev,_:=session["Evidence"].([]infermedica.Evidence)
newev=append(newev, infermedica.Evidence{ID:session["questioninfo"].(string),ChoiceID:infermedica.EvidenceChoiceIDPresent})  
session["Evidence"]=newev
return nil
 case "no": newev,_:=session["Evidence"].([]infermedica.Evidence)
newev=append(newev, infermedica.Evidence{ID:session["questioninfo"].(string),ChoiceID:infermedica.EvidenceChoiceIDAbsent}) 
session["Evidence"]=newev
return nil
 case "unknown": newev,_:=session["Evidence"].([]infermedica.Evidence)
newev=append(newev, infermedica.Evidence{ID:session["questioninfo"].(string),ChoiceID:infermedica.EvidenceChoiceIDUnknown}) 
session["Evidence"]=newev 
return nil
default:return errors.New("not a valid choice  please choose yes , no or unknown") 

}

}

func Getsymptoms(w http.ResponseWriter, r *http.Request){

       
    symptoms, err := app.Symptoms()
if err != nil {
 fmt.Printf("Could not fetch symptoms: %v", err)
}
//c := appengine.NewContext(r)
   json.NewEncoder(w).Encode(symptoms)

}

func Parseei (texti string,session Session)(error){
 var evarray [] infermedica.Evidence
 ev, evFound := session["Evidence"]
    if !evFound {

  }else{

evarray=ev.([]infermedica.Evidence)

  }

  var parsein infermedica.ParseReq
  parsein.Text=texti
     ParseOutput,_:= app.Parse(parsein)
  
     for i := 0; i < len(ParseOutput.Mentions); i++ { 
   
    
        
        if ParseOutput.Mentions[i].ChoiceID=="present"{
         evarray=append(evarray, infermedica.Evidence{ID:ParseOutput.Mentions[i].ID,ChoiceID:infermedica.EvidenceChoiceIDPresent})
        }


        if ParseOutput.Mentions[i].ChoiceID=="absent"{
             evarray=append(evarray, infermedica.Evidence{ID:ParseOutput.Mentions[i].ID,ChoiceID:infermedica.EvidenceChoiceIDAbsent})
        }
    if ParseOutput.Mentions[i].ChoiceID=="unknown"{
        evarray=append(evarray, infermedica.Evidence{ID:ParseOutput.Mentions[i].ID,ChoiceID:infermedica.EvidenceChoiceIDUnknown})
        }
    
    }

    if len(ParseOutput.Mentions)==0 {
    return errors.New("empty") 
}


    session["Evidence"]=evarray
     return nil
   
}

func diagnosis(DOutput  DiagnosisReq) (*infermedica.DiagnosisRes){

    //app := infermedica.NewApp("a3e97cd7", "af637361b0854d49d87b6af528c086dc", "")
       
      
 
  
DiagOutput,_ := Diagnosis(DOutput)
  
 
 return DiagOutput
  
     //json.NewEncoder(w).Encode(&ev)
    //json.NewEncoder(w).Encode(&firstinput)
}


// handle Handles /
func handle(w http.ResponseWriter, r *http.Request) {
    body :=
        "<!DOCTYPE html><html><head><title>Chatbot</title></head><body><pre style=\"font-family: monospace;\">\n" +
            "Available Routes:\n\n" +
            "  GET  /welcome -> handleWelcome\n" +
            "  POST /chat    -> handleChat\n" +
            "  GET  /        -> handle        (current)\n" +
            "</pre></body></html>"
    w.Header().Add("Content-Type", "text/html")
    fmt.Fprintln(w, body)
}



func handleChat(w http.ResponseWriter, r *http.Request) {
   // Make sure only POST requests are handled
    // Make sure a UUID exists in the Authorization header
    if r.Method != http.MethodPost {
        http.Error(w, "Only POST requests are allowed.", http.StatusMethodNotAllowed)
        return
    }
    uuid := r.Header.Get("Authorization")
    if uuid == "" {
        http.Error(w, "Missing or empty Authorization header.", http.StatusUnauthorized)
        return
    }

    // Make sure a session exists for the extracted UUID
    session, sessionFound := sessions[uuid]
    if !sessionFound {
        http.Error(w, fmt.Sprintf("No session found for: %v.", uuid), http.StatusUnauthorized)
        return
    }

    // Parse the JSON string in the body of the request
    data := JSON{}
    if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
        http.Error(w, fmt.Sprintf("Couldn't decode JSON: %v.", err), http.StatusBadRequest)
        return
    }
    defer r.Body.Close()

    // Make sure a message key is defined in the body of the request
    _, messageFound := data["message"]
    if !messageFound {
        http.Error(w, "Missing message key in body.", http.StatusBadRequest)
        return
    }

    // Process the received message
    message, err := processor(session, data["message"].(string))
    if err != nil {
        http.Error(w, err.Error(), 422 /* http.StatusUnprocessableEntity */)
        return
    }

   
    // Here we will call the diagones mehod and put the question in the message filed 
    writeJSON(w, JSON{
        "message":message,
    })

}



func handleWelcome(w http.ResponseWriter, r *http.Request) {
    u1 := uuid.NewV4()
    l:= u1.String()
 
sessions[l] = Session{}
     writeJSON(w, JSON{
        "uuid":    u1,
        "message": WelcomeMessage,
    })

 
   //session[l]=u1
}
func withLog(fn http.HandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        c := httptest.NewRecorder()
        fn(c, r)
        log.Printf("[%d] %-4s %s\n", c.Code, r.Method, r.URL.Path)

        for k, v := range c.HeaderMap {
            w.Header()[k] = v
        }
        w.WriteHeader(c.Code)
        c.Body.WriteTo(w)
    }
}


func writeJSON(w http.ResponseWriter, data JSON) {
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(data)
}

func  Diagnosis(dr DiagnosisReq) (*infermedica.DiagnosisRes, error) {
    if !dr.Sex.IsValid() {
        return nil, errors.New("Unexpected value for Sex")
    }
    b := new(bytes.Buffer)
    err1 := json.NewEncoder(b).Encode(dr)
    if err1 != nil {
        return nil, err1
    }
//app key
    fmt.Println("l1")
    req, err := http.NewRequest("POST", "https://api.infermedica.com/v2/diagnosis", b)
    if err != nil {
        fmt.Println("error")
        return nil, err
    }
    //a3e97cd7", "af637361b0854d49d87b6af528c086dc"
var A string ="a3e97cd7"
var B string="af637361b0854d49d87b6af528c086dc"
    fmt.Println(req)
    //req.Header.Add("Strict-Transport-Security", "max-age=63072000; includeSubDomains")
//req.Header.Add("Access-Control-Allow-Origin", "*")
    
    req.Header.Add("App-Id",A)
    req.Header.Add("App-Key",B)
    req.Header.Add("Content-Type", "application/json")
    tr := &http.Transport{
        TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
    }
    client := &http.Client{Transport: tr}
    res, err3 := client.Do(req)

    if err3 != nil {
        fmt.Println(err3)
        return nil, err3
    }
    fmt.Println(res)
    r := infermedica.DiagnosisRes{}
    err = json.NewDecoder(res.Body).Decode(&r)
     fmt.Println(r)
    if err != nil {
        return nil, err
    }
    return &r, nil
}


func main() {

ProcessFunc(sampleProcessor)

    // Use the PORT environment variable
    port := os.Getenv("PORT")
    // Default to 3000 if no PORT environment variable was defined
    if port == "" {
        port = "3000"
    }

    // Start the server
    fmt.Printf("Listening on port %s...\n", port)
    log.Fatalln(Engage(":" + port))









//     router := mux.NewRouter()
// chatbot.ProcessFunc(sampleProcessor)
//     router.HandleFunc("/", withLog(handle)).Methods("GET")
//     router.HandleFunc("/symptoms", Getsymptoms).Methods("GET")
//      router.HandleFunc("/welcome", withLog(handleWelcome)).Methods("GET")
// router.HandleFunc("/chat", withLog(handleChat)).Methods("POST")
//     log.Fatal(http.ListenAndServe(":8000", router))
}


func Engage(addr string) error {
    // HandleFuncs
    mux := http.NewServeMux()
    mux.HandleFunc("/welcome", withLog(handleWelcome))
    mux.HandleFunc("/chat", withLog(handleChat))
    mux.HandleFunc("/", withLog(handle))

    // Start the server
    return http.ListenAndServe(addr, cors.CORS(mux))
}


//a3c6fd82-048e-476f-a79e-7d2a5b1c374a