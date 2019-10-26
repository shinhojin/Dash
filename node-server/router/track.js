var express = require('express')
var router = express.Router()
var path = require('path')
var execSync = require('child_process').execSync;

router.post("/", function(req, res, next){
    result = {
        "tree" : doQuery(req.body.comment)
    }
    if(result == null){
        console.log("error")
    }
    //console.log(result)
    //console.log("=======")
    //console.log(result.SubId[0])
    res.render('track', result)
})

var doQuery = function(param){
    //console.log("CALL doQueary: " + param)
    var result = execSync('docker exec cli /bin/bash -c "./scripts/query.sh 0 1 "'+ param)
    var pattern=/\{.+?\}/g;
    var match = pattern.exec(result)
    var match2 = match[0].replace(/\\/gi,"")
    var data = JSON.parse(match2)
    var subs = data.SubId
    var node = {
            "nodeName" : data.OrderId,
            "name"     : data.ProductId,
            "type"     : data.Id.replace(/[0-9]/g, ""),
            "code"     : data.Detail,
            "label"    : data.OrderId,
            "version"  : "v1.0",
            "link" : {
				"name" : data.ProductId,
				"nodeName" : data.OrderId,
				"direction" : "SYNC"
			},
    }
    //console.log("THIS IS FIRST DATA")
    //console.log(data)
    if (subs == 'NULL'){
        //data.SubId = 'null'
        //console.log("sub is null")
        node["children"] = []
    }else{
        var pattern2 = /\%.+?\%/g;
        var subArray0 = subs.match(pattern2)
        //console.log("THIS IS SUB QUEARY")
        //console.log(subArray0)
        var retn = []
        for(i in subArray0 ){
            //console.log("=======RESURCIVE CALL=======")
            tmp = subArray0[i].replace(/\%/gi,"")
            retn[i] = doQuery(tmp)
        }
        console.log("THIS IS SUB QUEARY RESULT")
        console.log(retn)
        //data.SubId = retn
        node["children"]=retn
    }  
    //console.log("RESULT OF :"+param)
    //console.log(data)
    //console.log("END of doQuery"+param)
    return node
}

module.exports = router