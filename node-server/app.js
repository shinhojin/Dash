var express = require('express');
var app = express();
var path = require('path')
var bodyParser = require('body-parser')
var cookieParser = require('cookie-parser')
var session = require('express-session')


app.use(bodyParser.urlencoded({ extended: false }))
app.use(bodyParser.json())
app.use(express.static('public'));

app.use(cookieParser())
app.use(session({
	secret: '@#@$MYSIGN#@$#$',
	resave: false,
	saveUninitialized: true,
	cookie: {
		maxAge: 1000 * 60 * 60 * 24 * 365	// 쿠키 유효기간 1년
	},
}));

app.set('views', path.join(__dirname, 'view'))
app.set('view engine', 'ejs');

// test pages
var index = require('./router/index')
app.use('/index', index)

var track = require('./router/track')
app.use('/track', track)

var make = require('./router/make')
app.use('/make', make)

var move = require('./router/move')
app.use('/move', move)

var use = require('./router/use')
app.use('/use', use)

var hist = require('./router/history')
app.use('/history', hist)

app.use(function (err, req, res, next) {
	console.error(err);
	res.end("<h1>ERROR!</h1>")
});


app.listen(3000, function () {
    console.log('port 3000 opened');
});
  