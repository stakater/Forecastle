# ![Forecastle](assets/web/forecastle-round-100px.png) Forecastle

[![Get started with Stakater](https://stakater.github.io/README/stakater-github-banner.png)](http://stakater.com/?utm_source=IngressMonitorController&utm_medium=github)

## Problem(s)

- We would like to have a central place where we can easily look for and access our applications running on Kubernetes.
- We would like to have a tool which can dynamically discover and list the apps running on Kubernetes.
- A launchpad to access developer tools e.g. Jenkins, Nexus, Kibana, Grafana, etc.

## Solution

Forecastle gives you access to a control panel where you can see your running applications and access them on Kubernetes.

![Screenshot](assets/forecastle.png)

## Deploying to Kubernetes

You can deploy Forecastle both using vanilla k8s manifests or helm charts.

### Vanilla Manifests

#### Step 1: You can apply vanilla manifests by running the following command

```bash
kubectl apply -f https://raw.githubusercontent.com/stakater/Forecastle/master/deployments/kubernetes/forecastle.yaml
```

#### Step 2: Update configmap

In the Forecastle configmap modify the `namespaceSelector` key with a list of namespaces which you want Forecastle to watch. Refer to [this](#namespaceselector) for instructions.

And enjoy!

### Helm Charts

If you configured `helm` on your cluster, you can deploy Forecastle via helm chart located under `deployments/kubernetes/chart/Forecastle` folder.

## Configuration

### Ingresses

Forecastle looks for a specific annotations on ingresses.

- Add the following annotations to your ingresses in order to be discovered by forecastle:

| Annotation                                   | Description                                                                                                                                                 | Required |
| -------------------------------------------- | ----------------------------------------------------------------------------------------------------------------------------------------------------------- | -------- |
| `forecastle.stakater.com/expose`             | Add this with value `true` to the ingress of the app you want to show in Forecastle                                                                         | `true`   |
| `forecastle.stakater.com/icon`               | Icon/Image URL of the application; An icons/logos/images collection repo [Icons](https://github.com/stakater/ForecastleIcons)                               | `false`  |
| `forecastle.stakater.com/appName`            | A custom name for your application. Use if you don't want to use name of the ingress                                                                        | `false`  |
| `forecastle.stakater.com/group`              | A custom group name. Use if you want the application to show in a different group than the namespace it is running in                                       | `false`  |
| `forecastle.stakater.com/instance`           | A comma separated list of name/s of the forecastle instance/s where you want this application to appear. Use when you have multiple forecastle dashboards   | `false`  |
| `forecastle.stakater.com/url`                | A URL for the forecastle app (This will override the ingress URL). It MUST begin with a scheme i.e., `http://` or `https://`                                | `false`  |
| `forecastle.stakater.com/properties`         | A comma separate list of `key:value` pairs for the properties. This will appear as an expandable list for the app                                             | `false`  |
| `forecastle.stakater.com/network-restricted` | Specify whether the app is network restricted or not (true or false)                                                                                        | `false`  |

### Forecastle

Forecastle supports the following configuration options that can be modified by either ConfigMap or `values.yaml` if you are using helm

|       Field       |                                                Description                                                 |         Default         | Type              |
| :---------------: | :--------------------------------------------------------------------------------------------------------: | :---------------------: | ----------------- |
| namespaceSelector | A fine grained namespace selector which uses a combination of hardcoded namespaces well as label selectors |        any: true        | NamespaceSelector |
| headerBackground  |                         Background color of the header (Specified in the CSS way)                          |          null           | string            |
| headerForeground  |                         Foreground color of the header (Specified in the CSS way)                          |          null           | string            |
|       title       |                                     Title for the forecastle dashboard                                     | "Forecastle - Stakater" | string            |
|   instanceName    |                                      Name of the forecastle instance                                       |           ""            | string            |
|    customApps     |                A list of custom apps that you would like to add to the forecastle instance                 |           {}            | []CustomApp       |
|    crdEnabled     |                                  Enables or disables `ForecastleApp` CRD                                   |          true           | bool              |

#### NamespaceSelector

It is a selector for selecting namespaces either selecting all namespaces or a list of namespaces, or filtering namespaces through labels.

|     Field     |                                          Description                                          | Default | Type                                                                                         |
| :-----------: | :-------------------------------------------------------------------------------------------: | :-----: | -------------------------------------------------------------------------------------------- |
|      any      | Boolean describing whether all namespaces are selected in contrast to a list restricting them |  false  | bool                                                                                         |
| labelSelector |                Filter namespaces based on kubernetes metav1.LabelSelector type                |  null   | [metav1.LabelSelector](https://godoc.org/k8s.io/apimachinery/pkg/apis/meta/v1#LabelSelector) |
|  matchNames   |                                    List of namespace names                                    |  null   | []string                                                                                     |

*Note:* If you specify both `labelSelector` and `matchNames`, forecastle will take a union of all namespaces matched and use them.

#### Custom Apps

If you want to add any apps that are not exposed through ingresses or are external to the cluster, you can use the custom apps feature. You can pass an array of custom apps inside the config.

| Field             | Description                               | Type              |
| ----------------- | ----------------------------------------- | ----------------- |
| name              | Name of the custom app                    | String            |
| icon              | URL of the icon for the custom app        | String            |
| url               | URL of the custom app                     | String            |
| group             | Group for the custom app                  | String            |
| properties        | Additional Properties of the app as a map | map[string]string |
| networkRestricted | Whether app is network restricted or not  | bool              |

#### ForecastleApp CRD

You can now create custom resources to add apps to forecastle dynamically. This decouples the application configuration from Ingresses as well as forecastle config. You can create the custom resource `ForecastleApp` like the following:

```yaml
apiVersion: forecastle.stakater.com/v1alpha1
kind: ForecastleApp
metadata:
  name: app-name
spec:
  name: My Awesome App
  group: dev
  icon: https://icon-url
  url: http://app-url
  networkRestricted: "false"
  properties:
    Version: 1.0
  instance: "" # Optional
```

##### Automatically discover URL's from Kubernetes Resources

Forecastle supports discovering URL's ForecastleApp CRD from the following resources:

- Ingress

The above type of resource that you want to discover URL from **MUST** exist in the same namespace as `ForecastleApp` CR. Then you can add the following to the CR:

```yaml
apiVersion: forecastle.stakater.com/v1alpha1
kind: ForecastleApp
metadata:
  name: app-name
spec:
  name: My Awesome App
  group: dev
  icon: https://icon-url
  urlFrom: # This is new
    ingressRef:
      name: my-app-ingress
```

The above CR will be picked up by forecastle and it will generate the App in the UI. This lets you bundle this custom resource with the app's helm chart which will make it a part of the deployment process.

*Note:* You have to enable CRD feature first if you have disabled it. You can do that by applying the CRD and specifying `crdEnabled: true` in forecastle config. If you're using the helm chart then CRDs are installed with the chart.

#### Example Config

An example of a config can be seen below

```yaml
namespaceSelector:
  labelSelector:
    matchLabels:
      component: redis
    matchExpressions:
      - {key: tier, operator: In, values: [cache]}
  matchNames:
  - test
title:
headerBackground:
headerForeground: "#ffffff"
instanceName: "Hello"
crdEnabled: false
customApps:
- name: Hello
  icon: http://hello
  url: http://helloicon
  group: Test
  properties:
    Version: 1.0
```

## Features

- List apps found in all namespaces listed in the configmap
- Search apps
- Grouped apps per namespace
- Configurable header (Title and colors)
- Multiple instance support
- Provide Custom apps
- CRD `ForecastleApp` for adding custom apps
- Custom groups and URLs for the apps
- Details per app

## Running multiple instances of forecastle

You can run multiple instances of forecastle by just deploying them in a different namespace and provided a list of namespaces to look for ingresses. 

However, if you want flexibility over which applications to show in a specific instance regardless of the namespace, then you need to first configure forecastle instances to be a named instances. 
You can do that by setting `instanceName` in forecastle configuration. 

Once you have the named instances, you can add `forecastle.stakater.com/instance` annotation to your ingresses to control which application will show in which instance of forecastle. 

You can also specify multiple instances of forecastle for the same ingress so that it shows up in multiple dashboards. 
For example, you have 2 instances running named `dev-dashboard` and `prod-dashboard`. 
You can add this in the ingress's instance annotation `dev-dashboard,prod-dashboard` and the ingress will come up in both dashboards.

When using Helm make sure you set a unique `nameOverride` (default `forecastle`) in the values to prevent conflicts between global resources (`ClusterRole` & `ClusterRoleBinding`).

## Help

**Got a question?**
File a GitHub [issue](https://github.com/stakater/Forecastle/issues), or send us an [email](mailto:stakater@gmail.com).

### Talk to us on Slack

Join and talk to us on the #tools-imc channel for discussing Forecastle

[![Join Slack](https://stakater.github.io/README/stakater-join-slack-btn.png)](https://slack.stakater.com/)
[![Chat](https://stakater.github.io/README/stakater-chat-btn.png)](https://stakater-community.slack.com/messages/CBM7Q80KX)

## Contributing

### Bug Reports & Feature Requests

Please use the [issue tracker](https://github.com/stakater/Forecastle/issues) to report any bugs or file feature requests.

### Developing

PRs are welcome. In general, we follow the "fork-and-pull" Git workflow.

 1. **Fork** the repo on GitHub
 2. **Clone** the project to your own machine
 3. **Commit** changes to your own branch
 4. **Push** your work back up to your fork
 5. Submit a **Pull request** so that we can review your changes

NOTE: Be sure to merge the latest from "upstream" before making a pull request!

## Changelog

View our closed [Pull Requests](https://github.com/stakater/Forecastle/pulls?q=is%3Apr+is%3Aclosed).

## License

Apache2 © [Stakater](http://stakater.com)

## About

### Why name Forecastle

Forecastle is the section of the upper deck of a ship located at the bow forward of the foremast. This Forecastle will act as a control panel and show all your running applications on Kubernetes having a particular annotation.

`Forecastle` is maintained by [Stakater][website]. Like it? Please let us know at <hello@stakater.com>

See [our other projects][community]
or contact us in case of professional services and queries on <hello@stakater.com>

  [website]: http://stakater.com/
  [community]: https://www.stakater.com/projects-overview.html
