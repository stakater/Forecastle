const express = require('express')
const app = express()

app.use(express.static('public'))

app.listen(3000, () => console.log('Control Panel running on port 3000!'))