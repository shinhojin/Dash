var express = require('express')
var router = express.Router()
var path = require('path')
var execSync = require('child_process').execSync;

router.get("/", function(req, res, next){
    data = null
    res.render('history', data)
})

router.post("/", function(req, res, next){
    result = execSync('docker exec cli /bin/bash -c "./scripts/history.sh 0 1 "'+ req.body.orderId)
    if(result == null){
        console.log("error")
    }
    console.log(typeof(result))
    console.log(result)
    var pattern=/\#.+?\#/g;
    //var matchArray = result.match(pattern)

    //console.log(matchArray)
    data ={
        "terminal" : result
    }
    //console.log(result)
    //console.log("=======")
    //console.log(result.SubId[0])
    res.render('history', data)
})

module.exports = router