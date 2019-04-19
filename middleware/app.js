const createError = require(`http-errors`);
const express = require(`express`);
const path = require(`path`);
const cookieParser = require(`cookie-parser`);
const logger = require(`morgan`);

const indexRouter = require(`./routes/index`);
const usersRouter = require(`./routes/users`);
const statementsRouter = require(`./routes/statements`);
const diagnosisRouter = require(`./routes/diagnosis`);

const app = express();

// view engine setup
app.set(`views`, path.join(__dirname, `views`));
app.set(`view engine`, `ejs`);

app.use(logger(`dev`));
app.use(express.json());
app.use(express.urlencoded({ extended: false }));
app.use(cookieParser());
app.use(express.static(path.join(__dirname, `public`)));

// ***** Middleware to allow CORS access.
app.use((req, res, next) => {
  res.setHeader(`Access-Control-Allow-Origin`, `*`);
  res.setHeader(`Access-Control-Allow-Methods`, `GET, POST, OPTIONS, PUT, DELETE`);
  res.setHeader(`Access-Control-Allow-Headers`, `X-Requested-With,Authorization,X-PINGOTHER,Content-Type`);
  // Check if this is a preflight request. If so, send 200. Otherwise, pass it forward.
  if (req.method === `OPTIONS`) {
    //respond with 200
    res.sendStatus(200);
  } else {
    next();
  }
});

app.use(`/`, indexRouter);
app.use(`/users`, usersRouter);
app.use(`/statements`, statementsRouter);
app.use(`/diagnosis`, diagnosisRouter);

// catch 404 and forward to error handler
app.use(function(req, res, next) {
  next(createError(404));
});

// error handler
app.use(function(err, req, res) {
  // set locals, only providing error in development
  res.locals.message = err.message;
  res.locals.error = req.app.get(`env`) === `development` ? err : {};

  // render the error page
  res.status(err.status || 500);
  res.render(`error`);
});

module.exports = app;
