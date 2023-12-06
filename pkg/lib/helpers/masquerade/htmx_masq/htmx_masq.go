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
            <div class="bg-secondary rounded d-flex align-items-center justify-content-between p-4">
                    <i class="fa fa-info fa-3x text-primary"></i><h5>Basic settings</h5>
                    <div class="ms-3">
                    
                    <p class="mb-2">Environment: <h6 class="mb-0">{{.Env}}</h6></p>
                    <hr>
                    <p class="mb-2">LogPath: <h6 class="mb-0">{{.LogPath}}</h6></p>
                    <hr>
                    <p class="mb-2">ShutdownTimeOut: <h6 class="mb-0">{{.ShutdownTimeOut}}</h6></p>
                    
                    </div>
            </div>
        </div>

        <div class="col-sm-6 col-xl-3">
            <div class="bg-secondary rounded d-flex align-items-center justify-content-between p-4">
                <i class="fa fa-tree fa-3x text-primary"></i><h5>Core settings</h5>
                <div class="ms-3">

                    <p class="mb-2">Bucket size: <h6 class="mb-0">{{.CoreSettings.BucketSize}}</h6></p>

                </div>
            </div>
        </div>
    </div>
</div>

<div class="container-fluid pt-4 px-4">
    <div class="row g-4">
        <div class="col-sm-6 col-xl-3">
            <div class="bg-secondary rounded d-flex align-items-center justify-content-between p-4">
                <i class="fa fa-exchange-alt fa-3x text-primary"></i><h5>Web Socket</h5>
                <div class="ms-3">
                
                <p class="mb-2">Enable: <h6 class="mb-0">{{.WebSocketConnector.Enable}}</h6></p>
                <hr>
                <p class="mb-2">Address: <h6 class="mb-0">{{.WebSocketConnector.Address}}</h6></p>
                <hr>
                <p class="mb-2">Port: <h6 class="mb-0">{{.WebSocketConnector.Port}}</h6></p>
                <hr>
                <p class="mb-2">BufferSize - Read: <h6 class="mb-0">{{.WebSocketConnector.BufferSize.Read}}</h6></p>
                <hr>
                <p class="mb-2">BufferSize - Write: <h6 class="mb-0">{{.WebSocketConnector.BufferSize.Write}}</h6></p>
            
                </div>
            </div>
        </div>
    
        <div class="col-sm-6 col-xl-3">
            <div class="bg-secondary rounded d-flex align-items-center justify-content-between p-4">
                <i class="fa fa-server fa-3x text-primary"></i><h5>REST</h5>
                <div class="ms-3">
                
                <p class="mb-2">Enable: <h6 class="mb-0">{{.RestConnector.Enable}}</h6></p>
                <hr>
                <p class="mb-2">Address: <h6 class="mb-0">{{.RestConnector.Address}}</h6></p>
                <hr>
                <p class="mb-2">Port: <h6 class="mb-0">{{.RestConnector.Port}}</h6></p>
                
                </div>
            </div>
        </div>

        <div class="col-sm-6 col-xl-3">
            <div class="bg-secondary rounded d-flex align-items-center justify-content-between p-4">
                <i class="fa fa-shapes fa-3x text-primary"></i><h5>gRPC</h5>
                <div class="ms-3">
                
                <p class="mb-2">Enable: <h6 class="mb-0">{{.GrpcConnector.Enable}}</h6></p>
                <hr>
                <p class="mb-2">Address: <h6 class="mb-0">{{.GrpcConnector.Address}}</h6></p>
                <hr>
                <p class="mb-2">Port: <h6 class="mb-0">{{.GrpcConnector.Port}}</h6></p>
                
                </div>
            </div>
        </div>
    
        <div class="col-sm-6 col-xl-3">
            <div class="bg-secondary rounded d-flex align-items-center justify-content-between p-4">
                <i class="fa fa-window-restore fa-3x text-primary"></i><h5>Web Manager</h5>
                <div class="ms-3">
                
                <p class="mb-2">Enable: <h6 class="mb-0">{{.WebServer.Enable}}</h6></p>
                <hr>
                <p class="mb-2">Address: <h6 class="mb-0">{{.WebServer.Address}}</h6></p>
                <hr>
                <p class="mb-2">Port: <h6 class="mb-0">{{.WebServer.Port}}</h6></p>
                
                </div>
            </div>
        </div>
    
    </div>
</div>
`
