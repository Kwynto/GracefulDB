package htmx_masq

var Databases string = `
<div class="container-fluid pt-4 px-4">
    <div class="bg-secondary text-center rounded p-4">
	    <h4>Databases</h4>
    </div>
</div>

<div class="container-fluid pt-4 px-4">
    <div class="bg-secondary rounded p-4" id="dbAnswerOut">
    </div>
</div>

<!-- Footer Start -->
<div class="container-fluid pt-4 px-4">
    <div class="bg-secondary rounded-top p-4">
        <form hx-post="/hx/databases/request" hx-target="#dbAnswerOut" hx-swap="beforeend" hx-trigger="submit">
            <div class="d-flex mb-2">
                <input type="text" class="form-control bg-dark border-0" name="request" id="request-input" placeholder="Query">
                <input type="submit" class="btn btn-success ms-2" value="Send">
            </div>
        </form>
    </div>
</div>
<!-- Footer End -->
`

var DatabaseRequestAnswer string = `
    <p class="text-white"><b class="text-primary">Request [{{.TimeR}}] from {{.From}} &gt;</b> {{.Request}}</p>
    <p class="text-info"><b class="text-success">Answer [{{.TimeA}}] &gt;</b> {{.Answer}}</p>
`
