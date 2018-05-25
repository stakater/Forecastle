const express = require('express');
const app = express();
const bodyParser = require('body-parser');
const execSync = require('child_process').execSync;


app.use(express.static('public'));
app.use(bodyParser.urlencoded({ extended: true }));
app.use(bodyParser.json());

var router = express.Router(); 
router.get('/apps', function(req, res) {
    fs = require('fs');
    fs.readFile('/etc/cp-config/namespaces.conf', 'utf8', function (err,data) {
        if (err) {
            return console.log(err);
        }     
        var cmd = 'stk list ingresses --namespaces '+data+' --file /app/public/apps.json';
        execSync(cmd)
        res.send("")
    });
});

app.use('/api', router);

app.listen(3000, () => console.log('Control Panel running on port 3000!'))