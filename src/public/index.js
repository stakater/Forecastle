$(document).ready(function() {
    $.getJSON('apps.json', populateAppsListFromJson);
});

function populateAppsListFromJson(data) {
    $.each(data, renderApp);
}

function renderApp(key, app) {
    var appTemplate = $("#app-template").html();
    var html = Mustache.render(appTemplate, app)
    $(".appsList").append(html);
}