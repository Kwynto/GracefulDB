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
        <p>In this section, you can manage DBMS users.</p>
    </div>
</div>

<div class="container-fluid pt-4 px-4">
    <div class="row g-4">
        <div class="col-12">
            <div class="bg-secondary rounded h-100 p-4">
                <button type="button" class="btn btn-success" data-bs-toggle="modal" data-bs-target="#createModal" hx-get="/hx/accounts/create_load_form" hx-target="#createModalSpace"><i class="fa fa-plus-square"></i> Create a user</button><br>
                <table class="table table-hover">
                    <thead>
                        <tr>
                            <th scope="col">#</th>
                            <th scope="col">Login</th>
                            <th scope="col">Role</th>
                            <th scope="col">Description</th>
                            <th scope="col">Control</th>
                        </tr>
                    </thead>
                    <tbody>
                    {{ range $ind, $data := . }}
                        <tr>
                            <th scope="row"> {{ $ind }} </th>
                            <td> {{ $data.Login }} </td>
                            <td> {{ $data.Role }} </td>
                            <td> {{ $data.Description }} </td>
                            <td>
                                <div class="btn-group" role="group">
                                    <button type="button" hx-get="/hx/accounts/edit_form?user={{$data.Login}}" hx-target="#idMainUnit" class="btn btn-sm btn-info"><i class="fa fa-edit"></i> Edit</button>
                                    {{ if $data.Superuser }}
                                    <button type="button" class="btn btn-sm btn-warning" disabled><i class="fa fa-ban"></i> Block</button>
                                    <button type="button" class="btn btn-sm btn-danger" disabled><i class="fa fa-trash-alt"></i> Remove</button>
                                    {{ else }}
                                    <button type="button" hx-get="/hx/accounts/ban_form?user={{$data.Login}}" hx-target="#idMainUnit" class="btn btn-sm btn-warning"><i class="fa fa-ban"></i> Block</button>
                                    <button type="button" hx-get="/hx/accounts/del_form?user={{$data.Login}}" hx-target="#idMainUnit" class="btn btn-sm btn-danger"><i class="fa fa-trash-alt"></i> Remove</button>
                                    {{ end }}
                                </div> 
                            </td>
                        </tr>
                    {{ end }}
                    </tbody>
                </table>
            </div>
        </div>
    </div>
</div>

<div class="modal fade" id="createModal" tabindex="-1" aria-labelledby="createModalLabel" aria-hidden="true">
    <div class="modal-dialog modal-dialog-centered modal-dialog-scrollable">
        <div class="modal-content bg-light" id="createModalSpace">
            <div class="modal-header">
                <h1 class="modal-title fs-5" id="createModalLabel">Create user</h1>
                <button type="button" class="btn-close" data-bs-dismiss="modal" aria-label="Close"></button>
            </div>
            <div class="modal-body text-dark">
                <form id="create-user-form" hx-post="/hx/accounts/create_ok" hx-target="#createModalSpace" hx-trigger="submit">
                    <div class="mb-3">
                        <label for="login-input" class="col-form-label">Login:</label>
                        <input type="text" class="form-control" name="login" id="login-input">
                    </div>
                    <div class="mb-3">
                        <label for="password-input" class="col-form-label">Password:</label>
                        <input type="password" class="form-control" name="password" id="password-input">
                    </div>
                    <div class="mb-3">
                        <label for="desc-input" class="col-form-label">Description:</label>
                        <input type="text" class="form-control" name="desc" id="desc-input">
                    </div>
                </form>
            </div>
            <div class="modal-footer">
                <button type="button" class="btn btn-secondary" data-bs-dismiss="modal">Cancel</button>
                <button type="submit" form="create-user-form" class="btn btn-primary">Create</button>
            </div>
        </div>
    </div>
</div>
`

var AccountCreateFormLoad string = `
<div class="modal-header">
    <h1 class="modal-title fs-5" id="createModalLabel">Create user</h1>
    <button type="button" class="btn-close" data-bs-dismiss="modal" aria-label="Close"></button>
</div>
<div class="modal-body text-dark">
    <form id="create-user-form" hx-post="/hx/accounts/create_ok" hx-target="#createModalSpace" hx-trigger="submit">
        <div class="mb-3">
            <label for="login-input" class="col-form-label">Login:</label>
            <input type="text" class="form-control" name="login" id="login-input">
        </div>
        <div class="mb-3">
            <label for="password-input" class="col-form-label">Password:</label>
            <input type="password" class="form-control" name="password" id="password-input">
        </div>
        <div class="mb-3">
            <label for="desc-input" class="col-form-label">Description:</label>
            <input type="text" class="form-control" name="desc" id="desc-input">
        </div>
    </form>
</div>
<div class="modal-footer">
    <button type="button" class="btn btn-secondary" data-bs-dismiss="modal">Cancel</button>
    <button type="submit" form="create-user-form" class="btn btn-primary">Create</button>
</div>
`

var AccountCreateFormOk string = `
<div class="modal-header">
    <h1 class="modal-title fs-5" id="createModalLabel">Create user</h1>
    <button type="button" class="btn-close" data-bs-dismiss="modal" aria-label="Close"></button>
</div>
<div class="modal-body text-dark">
    Congratulations! The {{ .Login }} user has been created.<br>
</div>
<div class="modal-footer">
</div>
`

var AccountCreateFormError string = `
<div class="modal-header">
    <h1 class="modal-title fs-5" id="createModalLabel">Create user</h1>
    <button type="button" class="btn-close" data-bs-dismiss="modal" aria-label="Close"></button>
</div>
<div class="modal-body text-dark">
    User creation error.<br>
    The {{ .Login }} user cannot be created.<br>
</div>
<div class="modal-footer">
</div>
`

var AccountEditForm string = `
<div class="container-fluid pt-4 px-4">
    <div class="bg-secondary text-center rounded p-4">
	    <h4>Accounts</h4>
        <p>In this section, you can manage DBMS users.</p>
    </div>
</div>

<div class="container-fluid pt-4 px-4">
    <div class="bg-secondary text-center rounded p-4">
        Тут будет форма редактирования пользователя.
    </div>
</div>
`

var AccountBanForm string = `
<div class="container-fluid pt-4 px-4">
    <div class="bg-secondary text-center rounded p-4">
	    <h4>Accounts</h4>
        <p>In this section, you can manage DBMS users.</p>
    </div>
</div>

<div class="container-fluid pt-4 px-4">
    <div class="bg-secondary text-center rounded p-4">
        Тут будет форма блокировки пользователя.
    </div>
</div>
`

var AccountDelForm string = `
<div class="container-fluid pt-4 px-4">
    <div class="bg-secondary text-center rounded p-4">
	    <h4>Accounts</h4>
        <p>In this section, you can manage DBMS users.</p>
    </div>
</div>

<div class="container-fluid pt-4 px-4">
    <div class="bg-secondary text-center rounded p-4">
        Тут будет форма удаления пользователя.
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
                                    <input class="form-check-input" type="checkbox" role="switch" id="idWSCSwitch" hx-get="/hx/settings/wsc_change_sw" hx-target="#idMainUnit" hx-trigger="click delay:1s" checked>
                                {{else}}
                                    <input class="form-check-input" type="checkbox" role="switch" id="idWSCSwitch" hx-get="/hx/settings/wsc_change_sw" hx-target="#idMainUnit" hx-trigger="click delay:1s">
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
                                    <input class="form-check-input" type="checkbox" role="switch" id="idRestSwitch" hx-get="/hx/settings/rest_change_sw" hx-target="#idMainUnit" hx-trigger="click delay:1s" checked>
                                {{else}}
                                    <input class="form-check-input" type="checkbox" role="switch" id="idRestSwitch" hx-get="/hx/settings/rest_change_sw" hx-target="#idMainUnit" hx-trigger="click delay:1s">
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
                                    <input class="form-check-input" type="checkbox" role="switch" id="idGrpcSwitch" hx-get="/hx/settings/grpc_change_sw" hx-target="#idMainUnit" hx-trigger="click delay:1s" checked>
                                {{else}}
                                    <input class="form-check-input" type="checkbox" role="switch" id="idGrpcSwitch" hx-get="/hx/settings/grpc_change_sw" hx-target="#idMainUnit" hx-trigger="click delay:1s">
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
                                    <input class="form-check-input" type="checkbox" role="switch" id="idWebSwitch" hx-get="/hx/settings/web_change_sw" hx-target="#idMainUnit" hx-trigger="click delay:1s" disabled checked>
                                {{else}}
                                    <input class="form-check-input" type="checkbox" role="switch" id="idWebSwitch" hx-get="/hx/settings/web_change_sw" hx-target="#idMainUnit" hx-trigger="click delay:1s" disabled>
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
