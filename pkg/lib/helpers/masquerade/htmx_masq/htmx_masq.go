package htmx_masq

var Default string = `
<div class="container-fluid pt-4 px-4">
    <div class="bg-secondary text-center rounded p-4">
		<i class="fa fa-exclamation-triangle"></i> Error: Bad request
    </div>
</div>
`

var Dashboard string = `
<div class="container-fluid pt-4 px-4">
    <div class="bg-secondary text-center rounded p-4">
	    <h4>Dashboard</h4>
    </div>
</div>
`

var Databases string = `
<div class="container-fluid pt-4 px-4">
    <div class="bg-secondary text-center rounded p-4">
	    <h4>Databases</h4>
    </div>
</div>
`

var Accounts string = `
<div class="container-fluid pt-4 px-4">
    <div class="bg-secondary text-center rounded p-4">
	    <h4>Accounts</h4>
    </div>
</div>
`

var Settings string = `
<div class="container-fluid pt-4 px-4">
    <div class="bg-secondary text-center rounded p-4">
	    <h4>Settings</h4>
        <p>The server configuration can be changed via the configuration file.<br>
        Here you can only switch components quickly.</p>
    </div>
</div>

<div class="container-fluid pt-4 px-4">
    <div class="row g-4">
        <div class="col-sm-6 col-xl-3">
            <div class="bg-secondary rounded p-4">
                <h5 class="mb-4">Basic settings</h5>
                <table class="table table-hover">
                    <tbody>
                        <tr>
                            <td><h6>Environment: </h6></td>
                            <td>{{.Env}}</td>
                        </tr>
                        <tr>
                            <td><h6>LogPath: </h6></td>
                            <td>{{.LogPath}}</td>
                        </tr>
                        <tr>
                            <td><h6>ShutdownTimeOut: </h6></td>
                            <td>{{.ShutdownTimeOut}}</td>
                        </tr>
                    </tbody>
                </table>
            </div>
        </div>

        <div class="col-sm-6 col-xl-3">
            <div class="bg-secondary rounded p-4">
                <h5 class="mb-4">Core settings</h5>
                <table class="table table-hover">
                    <tbody>
                        <tr>
                            <td><h6>Bucket size: </h6></td>
                            <td>{{.CoreSettings.BucketSize}}</td>
                        </tr>
                    </tbody>
                </table>
            </div>
        </div>
    </div>
</div>

<div class="container-fluid pt-4 px-4">
    <div class="row g-4">
        <div class="col-sm-6 col-xl-3">
            <div class="bg-secondary rounded p-4">
                <h5 class="mb-4">Web Socket Connector</h5>
                <table class="table table-hover">
                    <tbody>
                        <tr>
                            <td><h6>Enable: </h6></td>
                            <td>
                                <div class="form-check form-switch">
                                {{if .WebSocketConnector.Enable}}
                                    <input class="form-check-input" type="checkbox" role="switch" id="idWSCSwitch" checked>
                                {{else}}
                                    <input class="form-check-input" type="checkbox" role="switch" id="idWSCSwitch">
                                {{end}}
                                </div>
                            </td>
                        </tr>
                        <tr>
                            <td><h6>Address: </h6></td>
                            <td>{{.WebSocketConnector.Address}}</td>
                        </tr>
                        <tr>
                            <td><h6>Port: </h6></td>
                            <td>{{.WebSocketConnector.Port}}</td>
                        </tr>
                        <tr>
                            <td><h6>BufferSize - Read: </h6></td>
                            <td>{{.WebSocketConnector.BufferSize.Read}}</td>
                        </tr>
                        <tr>
                            <td><h6>BufferSize - Write: </h6></td>
                            <td>{{.WebSocketConnector.BufferSize.Write}}</td>
                        </tr>
                    </tbody>
                </table>
            </div>
        </div>

        <div class="col-sm-6 col-xl-3">
            <div class="bg-secondary rounded p-4">
                <h5 class="mb-4">REST Connector</h5>
                <table class="table table-hover">
                    <tbody>
                        <tr>
                            <td><h6>Enable: </h6></td>
                            <td>
                                <div class="form-check form-switch">
                                {{if .RestConnector.Enable}}
                                    <input class="form-check-input" type="checkbox" role="switch" id="idRestSwitch" checked>
                                {{else}}
                                    <input class="form-check-input" type="checkbox" role="switch" id="idRestSwitch">
                                {{end}}
                                </div>
                            </td>
                        </tr>
                        <tr>
                            <td><h6>Address: </h6></td>
                            <td>{{.RestConnector.Address}}</td>
                        </tr>
                        <tr>
                            <td><h6>Port: </h6></td>
                            <td>{{.RestConnector.Port}}</td>
                        </tr>
                    </tbody>
                </table>
            </div>
        </div>

        <div class="col-sm-6 col-xl-3">
            <div class="bg-secondary rounded p-4">
                <h5 class="mb-4">gRPC Connector</h5>
                <table class="table table-hover">
                    <tbody>
                        <tr>
                            <td><h6>Enable: </h6></td>
                            <td>
                                <div class="form-check form-switch">
                                {{if .GrpcConnector.Enable}}
                                    <input class="form-check-input" type="checkbox" role="switch" id="idGrpcSwitch" checked>
                                {{else}}
                                    <input class="form-check-input" type="checkbox" role="switch" id="idGrpcSwitch">
                                {{end}}
                                </div>
                            </td>
                        </tr>
                        <tr>
                            <td><h6>Address: </h6></td>
                            <td>{{.GrpcConnector.Address}}</td>
                        </tr>
                        <tr>
                            <td><h6>Port: </h6></td>
                            <td>{{.GrpcConnector.Port}}</td>
                        </tr>
                    </tbody>
                </table>
            </div>
        </div>

        <div class="col-sm-6 col-xl-3">
            <div class="bg-secondary rounded p-4">
                <h5 class="mb-4">GracefulDB Web Manager</h5>
                <table class="table table-hover">
                    <tbody>
                        <tr>
                            <td><h6>Enable: </h6></td>
                            <td>
                                <div class="form-check form-switch">
                                {{if .WebServer.Enable}}
                                    <input class="form-check-input" type="checkbox" role="switch" id="idWebSwitch" disabled checked>
                                {{else}}
                                    <input class="form-check-input" type="checkbox" role="switch" id="idWebSwitch" disabled>
                                {{end}}
                                </div>
                            </td>
                        </tr>
                        <tr>
                            <td><h6>Address: </h6></td>
                            <td>{{.WebServer.Address}}</td>
                        </tr>
                        <tr>
                            <td><h6>Port: </h6></td>
                            <td>{{.WebServer.Port}}</td>
                        </tr>
                    </tbody>
                </table>
            </div>
        </div>
    
    </div>
</div>
`
