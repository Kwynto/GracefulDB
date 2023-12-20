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
                <button type="button" class="btn btn-info" hx-get="/hx/nav/accounts" hx-target="#idMainUnit"><i class="fa fa-sync-alt"></i></button>
                <button type="button" class="btn btn-success" data-bs-toggle="modal" data-bs-target="#createModal" hx-get="/hx/accounts/create_load_form" hx-target="#createModalSpace"><i class="fa fa-plus-square"></i> Create a user</button><br>
                <table class="table table-hover">
                    <thead>
                        <tr>
                            <th scope="col">#</th>
                            <th scope="col">Login</th>
                            <th scope="col">Status</th>
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
                            <td> {{ $data.Status }} </td>
                            <td> {{ $data.Role }} </td>
                            <td> {{ $data.Description }} </td>
                            <td>
                                <div class="btn-group" role="group">
                                    {{ if $data.System }}
                                    <button type="button" class="btn btn-sm btn-success" disabled><i class="fa fa-edit"></i> Edit</button>
                                    {{ else }}
                                    <button type="button" class="btn btn-sm btn-success" data-bs-toggle="modal" data-bs-target="#editModal" hx-get="/hx/accounts/edit_load_form?user={{$data.Login}}" hx-target="#editModalSpace"><i class="fa fa-edit"></i> Edit</button>
                                    {{ end }}
                                    {{ if or $data.Superuser $data.System }}
                                    <button type="button" class="btn btn-sm btn-warning" style="width: 100px;" disabled><i class="fa fa-ban"></i> Block</button>
                                    <button type="button" class="btn btn-sm btn-danger" disabled><i class="fa fa-trash-alt"></i> Remove</button>
                                    {{ else if $data.Baned }}
                                    <button type="button" class="btn btn-sm btn-warning" style="width: 100px;" data-bs-toggle="modal" data-bs-target="#unbanModal" hx-get="/hx/accounts/unban_load_form?user={{$data.Login}}" hx-target="#unbanModalSpace"><i class="fa fa-undo"></i> UnBlock</button>
                                    <button type="button" class="btn btn-sm btn-danger" data-bs-toggle="modal" data-bs-target="#delModal" hx-get="/hx/accounts/del_load_form?user={{$data.Login}}" hx-target="#delModalSpace"><i class="fa fa-trash-alt"></i> Remove</button>
                                    {{ else }}
                                    <button type="button" class="btn btn-sm btn-warning" style="width: 100px;" data-bs-toggle="modal" data-bs-target="#banModal" hx-get="/hx/accounts/ban_load_form?user={{$data.Login}}" hx-target="#banModalSpace"><i class="fa fa-ban"></i> Block</button>
                                    <button type="button" class="btn btn-sm btn-danger" data-bs-toggle="modal" data-bs-target="#delModal" hx-get="/hx/accounts/del_load_form?user={{$data.Login}}" hx-target="#delModalSpace"><i class="fa fa-trash-alt"></i> Remove</button>
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
`

var AccountCreateFormOk string = `
<div class="modal-header" hx-get="/hx/nav/accounts" hx-trigger="load" hx-target="#idMainUnit">
    <h1 class="modal-title fs-5" id="createModalLabel">Create user</h1>
    <button type="button" class="btn-close" data-bs-dismiss="modal" aria-label="Close"></button>
</div>
<div class="modal-body text-dark">
    Congratulations! The <b>{{.Login}}</b> user has been created.<br>
</div>
<div class="modal-footer">
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

var AccountCreateFormError string = `
<div class="modal-header">
    <h1 class="modal-title fs-5" id="createModalLabel">Create user</h1>
    <button type="button" class="btn-close" data-bs-dismiss="modal" aria-label="Close"></button>
</div>
<div class="modal-body text-dark">
    User creation error.<br>
    The <b>{{.Login}}</b> user cannot be created.<br>
</div>
<div class="modal-footer">
</div>
`

var AccountEditFormOk string = `
<div class="modal-header" hx-get="/hx/nav/accounts" hx-trigger="load" hx-target="#idMainUnit">
    <h1 class="modal-title fs-5" id="editModalLabel">Edit user</h1>
    <button type="button" class="btn-close" data-bs-dismiss="modal" aria-label="Close"></button>
</div>
<div class="modal-body text-dark">
    Congratulations! The <b>{{.Login}}</b> user has been updated.<br>
</div>
<div class="modal-footer">
</div>
`

var AccountEditFormLoad string = `
<div class="modal-header">
    <h1 class="modal-title fs-5" id="editModalLabel">Edit user</h1>
    <button type="button" class="btn-close" data-bs-dismiss="modal" aria-label="Close"></button>
</div>
{{ if eq .Role 0 }} 
<div class="modal-body text-dark">
    SYSTEM USER (This user cannot be changed.) 
</div>
{{ else }}
<div class="modal-body text-dark">
    <form id="edit-user-form" hx-post="/hx/accounts/edit_ok" hx-target="#editModalSpace" hx-trigger="submit">
        <div class="mb-3">
            <label for="login-input" class="col-form-label">Login:</label> <b>{{.Login}}</b>
            <input type="hidden" class="form-control" name="login" id="login-input" value="{{.Login}}">
        </div>
        <div class="mb-3">
            <label for="password-input" class="col-form-label">New password:</label>
            <input type="password" class="form-control" name="password" id="password-input" value="">
        </div>
        <div class="mb-3">
            <label for="desc-input" class="col-form-label">Description:</label>
            <input type="text" class="form-control" name="desc" id="desc-input" value="{{.Description}}">
        </div>
        {{ if ne .Login "root" }} 
        <div class="mb-3">
            <label for="status-select" class="col-form-label">Status:</label>{{ if eq .Status 0 }} UNDEFINED (You must select the status) {{ end }}
            <select class="form-select form-select-sm mb-3" aria-label=".form-select-sm" id="status-select" name="status">
                {{ if eq .Status 1 }} <option value="1" selected>NEW</option> {{ else }} <option value="1">NEW</option> {{ end }}
                {{ if eq .Status 2 }} <option value="2" selected>ACTIVE</option> {{ else }} <option value="2">ACTIVE</option> {{ end }}
                {{ if eq .Status 3 }} <option value="3" selected>BANED</option> {{ else }} <option value="3">BANED</option> {{ end }}
            </select>
        </div>
        <div class="mb-3">
            <label for="role-select" class="col-form-label">Role:</label>
            <select class="form-select form-select-sm mb-3" aria-label=".form-select-sm" id="role-select" name="role">
                {{ if eq .Role 1 }} <option value="1" selected>ADMIN</option> {{ else }} <option value="1">ADMIN</option> {{ end }}
                {{ if eq .Role 2 }} <option value="2" selected>MANAGER</option> {{ else }} <option value="2">MANAGER</option> {{ end }}
                {{ if eq .Role 3 }} <option value="3" selected>ENGINEER</option> {{ else }} <option value="3">ENGINEER</option> {{ end }}
                {{ if eq .Role 4 }} <option value="4" selected>USER</option> {{ else }} <option value="4">USER</option> {{ end }}
            </select>
        </div>
        <div class="mb-3">
            <label for="rules-area" class="col-form-label">Rules:</label>
            <textarea class="form-control" id="rules-area" name="rules" style="height: 100px;">{{.Rules}}</textarea>
        </div>
        {{ else }}
            You cannot change the permissions of this user.
        {{ end }}
    </form>
</div>
<div class="modal-footer">
    <button type="button" class="btn btn-secondary" data-bs-dismiss="modal">Cancel</button>
    <button type="submit" form="edit-user-form" class="btn btn-primary">Save</button>
</div>
{{ end }}
`

var AccountEditFormError string = `
<div class="modal-header">
    <h1 class="modal-title fs-5" id="editModalLabel">Edit user</h1>
    <button type="button" class="btn-close" data-bs-dismiss="modal" aria-label="Close"></button>
</div>
<div class="modal-body text-dark">
    User update error.<br>
    The <b>{{.Login}}</b> user cannot be updated.<br>
</div>
<div class="modal-footer">
</div>
`

var AccountBanFormOk string = `
<div class="modal-header" hx-get="/hx/nav/accounts" hx-trigger="load" hx-target="#idMainUnit">
    <h1 class="modal-title fs-5" id="banModalLabel">Block user</h1>
    <button type="button" class="btn-close" data-bs-dismiss="modal" aria-label="Close"></button>
</div>
<div class="modal-body text-dark">
    Congratulations! The <b>{{.Login}}</b> user has been blocked.<br>
</div>
<div class="modal-footer">
</div>
`

var AccountBanFormLoad string = `
<div class="modal-header">
    <h1 class="modal-title fs-5" id="banModalLabel">Block user</h1>
    <button type="button" class="btn-close" data-bs-dismiss="modal" aria-label="Close"></button>
</div>
<div class="modal-body text-dark">
    <form id="ban-user-form" hx-post="/hx/accounts/ban_ok" hx-target="#banModalSpace" hx-trigger="submit">
        <div class="mb-3">
            <input type="hidden" class="form-control" name="login" id="login-input" value="{{.Login}}">
        </div>
    </form>
    Block the <b>{{.Login}}</b> user?
</div>
<div class="modal-footer">
    <button type="button" class="btn btn-secondary" data-bs-dismiss="modal">Cancel</button>
    <button type="submit" form="ban-user-form" class="btn btn-primary">Block</button>
</div>
`

var AccountBanFormError string = `
<div class="modal-header">
    <h1 class="modal-title fs-5" id="banModalLabel">Block user</h1>
    <button type="button" class="btn-close" data-bs-dismiss="modal" aria-label="Close"></button>
</div>
<div class="modal-body text-dark">
    Error blocking the user.<br>
    The <b>{{.Login}}</b> user cannot be blocked.<br>
</div>
<div class="modal-footer">
</div>
`

var AccountUnBanFormOk string = `
<div class="modal-header" hx-get="/hx/nav/accounts" hx-trigger="load" hx-target="#idMainUnit">
    <h1 class="modal-title fs-5" id="unbanModalLabel">UnBlock user</h1>
    <button type="button" class="btn-close" data-bs-dismiss="modal" aria-label="Close"></button>
</div>
<div class="modal-body text-dark">
    Congratulations! The <b>{{.Login}}</b> user has been unblocked.<br>
</div>
<div class="modal-footer">
</div>
`

var AccountUnBanFormLoad string = `
<div class="modal-header">
    <h1 class="modal-title fs-5" id="unbanModalLabel">UnBlock user</h1>
    <button type="button" class="btn-close" data-bs-dismiss="modal" aria-label="Close"></button>
</div>
<div class="modal-body text-dark">
    <form id="unban-user-form" hx-post="/hx/accounts/unban_ok" hx-target="#unbanModalSpace" hx-trigger="submit">
        <div class="mb-3">
            <input type="hidden" class="form-control" name="login" id="login-input" value="{{.Login}}">
        </div>
    </form>
    UnBlock the <b>{{.Login}}</b> user?
</div>
<div class="modal-footer">
    <button type="button" class="btn btn-secondary" data-bs-dismiss="modal">Cancel</button>
    <button type="submit" form="unban-user-form" class="btn btn-primary">UnBlock</button>
</div>
`

var AccountUnBanFormError string = `
<div class="modal-header">
    <h1 class="modal-title fs-5" id="unbanModalLabel">UnBlock user</h1>
    <button type="button" class="btn-close" data-bs-dismiss="modal" aria-label="Close"></button>
</div>
<div class="modal-body text-dark">
    Error unblocking the user.<br>
    The <b>{{.Login}}</b> user cannot be unblocked.<br>
</div>
<div class="modal-footer">
</div>
`

var AccountDelFormOk string = `
<div class="modal-header" hx-get="/hx/nav/accounts" hx-trigger="load" hx-target="#idMainUnit">
    <h1 class="modal-title fs-5" id="delModalLabel">Remove user</h1>
    <button type="button" class="btn-close" data-bs-dismiss="modal" aria-label="Close"></button>
</div>
<div class="modal-body text-dark">
    Congratulations! The <b>{{.Login}}</b> user has been deleted.<br>
</div>
<div class="modal-footer">
</div>
`

var AccountDelFormLoad string = `
<div class="modal-header">
    <h1 class="modal-title fs-5" id="delModalLabel">Remove user</h1>
    <button type="button" class="btn-close" data-bs-dismiss="modal" aria-label="Close"></button>
</div>
<div class="modal-body text-dark">
    <form id="del-user-form" hx-post="/hx/accounts/del_ok" hx-target="#delModalSpace" hx-trigger="submit">
        <div class="mb-3">
            <input type="hidden" class="form-control" name="login" id="login-input" value="{{.Login}}">
        </div>
    </form>
    Delete the <b>{{.Login}}</b> user?
</div>
<div class="modal-footer">
    <button type="button" class="btn btn-secondary" data-bs-dismiss="modal">Cancel</button>
    <button type="submit" form="del-user-form" class="btn btn-primary">Delete</button>
</div>
`

var AccountDelFormError string = `
<div class="modal-header">
    <h1 class="modal-title fs-5" id="delModalLabel">Remove user</h1>
    <button type="button" class="btn-close" data-bs-dismiss="modal" aria-label="Close"></button>
</div>
<div class="modal-body text-dark">
    Error deleting the user.<br>
    The <b>{{.Login}}</b> user cannot be deleted.<br>
</div>
<div class="modal-footer">
</div>
`
