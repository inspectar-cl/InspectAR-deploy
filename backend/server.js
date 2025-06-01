const express = require("express");
const cors = require("cors");
const morgan = require("morgan");
const dotenv = require("dotenv"); // No sé si lo utilizaremos, supongo que para más adelante del desarrollo sí (Vixo 31-05).
const bodyParser = require("body-parser");

const app = express();
app.use(morgan("dev"));
app.use(express.json());




// Start server
// const port = process.env.PORT; // || 3000;
const port = 3000
app.listen(port, () => {
  console.log(`Server listening on port ${port}`);
});
