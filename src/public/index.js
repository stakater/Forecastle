$(document).ready(function() {
    var fileName = "apps.json"
    if(FileExists(fileName)) {
        $.getJSON(fileName, populateAppsListFromJson);
    } else {
    $.ajax({
        url: "/api/apps",
        type: 'GET',
        success: function(res) {
            $.getJSON(fileName, populateAppsListFromJson);
        }
    });
    }
});

function populateAppsListFromJson(data) {
    $.each(data, renderApp);
}

function renderApp(key, app) {
    var appTemplate = $("#app-template").html();
    var html = Mustache.render(appTemplate, app)
    $(".appsList").append(html);
}

function FileExists(file)
{
    var http = new XMLHttpRequest();
    http.open('HEAD', file, false);
    http.send();
    return http.status!=404;
}