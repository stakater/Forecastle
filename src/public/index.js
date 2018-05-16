window.onload = function() {
    readTextFile("apps.json", function(data) {
        apps = JSON.parse(data)

        appsListUl = document.getElementById('appsList')

        for(var i = 0; i < apps.length; i++) {
            var li = document.createElement("li");
            var a = document.createElement("a");
            a.setAttribute('href', apps[i].URL);
            a.appendChild(document.createTextNode(apps[i].Name))
            li.appendChild(a);
            appsListUl.appendChild(li);
        }
    });
}

function readTextFile(file, callback)
{
    var rawFile = new XMLHttpRequest();
    rawFile.open("GET", file, false);
    rawFile.onreadystatechange = function ()
    {
        if(rawFile.readyState === 4)
        {
            if(rawFile.status === 200 || rawFile.status == 0)
            {
                var allText = rawFile.responseText;
                callback(allText)
            }
        }
    }
    rawFile.send(null);
}