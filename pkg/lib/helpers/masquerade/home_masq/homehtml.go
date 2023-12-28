package home_masq

var html1 string = `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <title>GracefulDB Manager</title>
    <meta content="width=device-width, initial-scale=1.0" name="viewport">

    <!-- Favicon -->
    <link href="./static/img/favicon.ico" rel="icon">

    <!-- Google Web Fonts -->
    <link rel="preconnect" href="https://fonts.googleapis.com">
    <link rel="preconnect" href="https://fonts.gstatic.com" crossorigin>
    <link href="https://fonts.googleapis.com/css2?family=Open+Sans:wght@400;600&family=Roboto:wght@500;700&display=swap" rel="stylesheet"> 
    
    <!-- Icon Font Stylesheet -->
    <link href="https://cdnjs.cloudflare.com/ajax/libs/font-awesome/5.10.0/css/all.min.css" rel="stylesheet">
    <link href="https://cdn.jsdelivr.net/npm/bootstrap-icons@1.4.1/font/bootstrap-icons.css" rel="stylesheet">

    <!-- Customized Bootstrap Stylesheet -->
	<style>
	`

var html2 string = `
	</style>

    <!-- Template Stylesheet -->
	<style>
	`

var html3 string = `
	</style>
</head>

<body hx-get="/hx/nav/dashboard" hx-trigger="load" hx-target="#idMainUnit">
    <div class="container-fluid position-relative d-flex p-0">
        <!-- Spinner Start -->
        <div id="spinner" class="show bg-dark position-fixed translate-middle w-100 vh-100 top-50 start-50 d-flex align-items-center justify-content-center">
            <div class="spinner-border text-primary" style="width: 3rem; height: 3rem;" role="status">
                <span class="sr-only">Loading...</span>
            </div>
        </div>
        <!-- Spinner End -->

        <!-- Sidebar Start -->
        <div class="sidebar pe-4 pb-3">
            <nav class="navbar bg-secondary navbar-dark">
                <a href="/" class="navbar-brand mx-4 mb-3">
                    <h3 class="text-primary">
                        <img src="./static/img/logo.svg" style="width: 50px; height: 50px;"> GracefulDB
                        <div class="tagline">Fast, Simple and Secure</div>
                    </h3>
                </a>
                <div class="d-flex align-items-center ms-4 mb-4">
                    <div class="position-relative">
                        <img class="rounded-circle" src="./static/img/user.png" alt="" style="width: 40px; height: 40px;">
                        <div class="bg-success rounded-circle border border-2 border-white position-absolute end-0 bottom-0 p-1"></div>
                    </div>
                    <div class="ms-3">
                        <h6 class="mb-0">{{ .Login }}</h6>
                        <span>{{ .Roles }}</span>
                    </div>
                </div>
                <div class="navbar-nav w-100">
                    <a hx-get="/hx/nav/dashboard" hx-target="#idMainUnit" class="nav-item nav-link"><i class="fa fa-tachometer-alt me-2"></i>Dashboard</a>
                    <a hx-get="/hx/nav/databases" hx-target="#idMainUnit" class="nav-item nav-link"><i class="fa fa-database me-2"></i>Databases</a>
                    <a hx-get="/hx/nav/accounts" hx-target="#idMainUnit" class="nav-item nav-link"><i class="fa fa-users me-2"></i>Accounts</a>
                    <a hx-get="/hx/nav/settings" hx-target="#idMainUnit" class="nav-item nav-link"><i class="fa fa-cogs me-2"></i>Settings</a>
                </div>
            </nav>
        </div>
        <!-- Sidebar End -->

        <!-- Content Start -->
        <div class="content">
            <!-- Navbar Start -->
            <nav class="navbar navbar-expand bg-secondary navbar-dark sticky-top px-4 py-0">
                <a href="/" class="navbar-brand d-flex d-lg-none me-4">
                    <h2 class="text-primary mb-0"><i class="fa fa-user-edit"></i></h2>
                </a>
                <a href="#" class="sidebar-toggler flex-shrink-0">
                    <i class="fa fa-bars"></i>
                </a>
                
                <div class="navbar-nav align-items-center ms-auto">
                    <div class="nav-item dropdown">
                        <a href="#" class="nav-link dropdown-toggle" data-bs-toggle="dropdown">
                            <img class="rounded-circle me-lg-2" src="./static/img/user.png" alt="" style="width: 40px; height: 40px;">
                            <span class="d-none d-lg-inline-flex">{{ .Login }}</span>
                        </a>
                        <div class="dropdown-menu dropdown-menu-end bg-secondary border-1 rounded-0 rounded-bottom m-0">
                            <a class="dropdown-item" hx-get="/hx/nav/dashboard" hx-target="#idMainUnit"><i class="fa fa-tachometer-alt me-2"></i>Dashboard</a>
                            <a class="dropdown-item" hx-get="/hx/nav/databases" hx-target="#idMainUnit"><i class="fa fa-database me-2"></i>Databases</a>
                            <a class="dropdown-item" hx-get="/hx/nav/accounts" hx-target="#idMainUnit"><i class="fa fa-users me-2"></i>Accounts</a>
                            <a class="dropdown-item" hx-get="/hx/nav/settings" hx-target="#idMainUnit"><i class="fa fa-cogs me-2"></i>Settings</a>
                            <hr>
                            <a class="dropdown-item" data-bs-toggle="modal" data-bs-target="#profileModal"><i class="fa fa-user me-2"></i>Your profile</a>
                            <a class="dropdown-item" data-bs-toggle="modal" data-bs-target="#logoutModal"><i class="fa fa-sign-out-alt me-2"></i>Log Out </a>
                        </div>
                    </div>
                </div>
            </nav>
            <!-- Navbar End -->

            <div id="idMainUnit">
            </div>

            <!-- Footer Start -->
            <div class="container-fluid pt-4 px-4">
                <div class="bg-secondary rounded-top p-4">
                    <div class="row">
                        <div class="col-12 col-sm-6 text-center text-sm-start">
                            &copy; <a target="_blank" href="http://gracefuldb.dev/">GracefulDB</a> 
                        </div>
                        <div class="col-12 col-sm-6 text-center text-sm-end">
                            <br>
                        </div>
                    </div>
                </div>
            </div>
            <!-- Footer End -->
        </div>
        <!-- Content End -->

        <!-- Back to Top -->
        <a href="#" class="btn btn-lg btn-primary btn-lg-square back-to-top"><i class="bi bi-arrow-up"></i></a>
    </div>

    <!-- Modal Profile -->
    <div class="modal fade" id="profileModal" tabindex="-1" aria-labelledby="profileModalLabel" aria-hidden="true">
        <div class="modal-dialog modal-dialog-centered modal-dialog-scrollable">
            <div class="modal-content bg-light" id="profileModalSpace">
                <div class="modal-header">
                    <h1 class="modal-title fs-5" id="profileModalLabel">Your profile</h1>
                    <button type="button" class="btn-close" data-bs-dismiss="modal" aria-label="Close"></button>
                </div>
                <div class="modal-body text-dark">
                    <form id="profile-user-form" hx-post="/hx/accounts/profile_ok" hx-target="#profileModalSpace" hx-trigger="submit">
                        <div class="mb-3">
                            <label for="login-input" class="col-form-label">Login:</label>
                            <input type="hidden" class="form-control" name="login" id="login-input" value="" disabled>
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
                    <button type="submit" form="profile-user-form" class="btn btn-primary">Save</button>
                </div>
            </div>
        </div>
    </div>

    <!-- Modal Logout -->
    <div class="modal fade" id="logoutModal" tabindex="-1" aria-labelledby="logoutModalLabel" aria-hidden="true">
        <div class="modal-dialog">
            <div class="modal-content bg-light">
                <div class="modal-header">
                    <h1 class="modal-title fs-5" id="logoutModalLabel">Log Out</h1>
                    <button type="button" class="btn-close" data-bs-dismiss="modal" aria-label="Close"></button>
                </div>
                <div class="modal-body text-dark">
                    Are you sure you wish to log out?
                </div>
                <div class="modal-footer">
                    <button type="button" class="btn btn-secondary" data-bs-dismiss="modal">No, I'll be back.</button>
                    <button type="button" class="btn btn-primary" hx-get="/hx/nav/logout">Yes, but I'm going to miss you.</button>
                </div>
            </div>
        </div>
    </div>

    <!-- Modals Accounts -->
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

    <div class="modal fade" id="editModal" tabindex="-1" aria-labelledby="editModalLabel" aria-hidden="true">
        <div class="modal-dialog modal-dialog-centered modal-dialog-scrollable">
            <div class="modal-content bg-light" id="editModalSpace">
                <div class="modal-header">
                    <h1 class="modal-title fs-5" id="editModalLabel">Edit user</h1>
                    <button type="button" class="btn-close" data-bs-dismiss="modal" aria-label="Close"></button>
                </div>
                <div class="modal-body text-dark">
                    <form id="edit-user-form" hx-post="/hx/accounts/edit_ok" hx-target="#editModalSpace" hx-trigger="submit">
                        <div class="mb-3">
                            <label for="login-input" class="col-form-label">Login:</label>
                            <input type="hidden" class="form-control" name="login" id="login-input" value="" disabled>
                        </div>
                        <div class="mb-3">
                            <label for="password-input" class="col-form-label">New password:</label>
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
                    <button type="submit" form="edit-user-form" class="btn btn-primary">Save</button>
                </div>
            </div>
        </div>
    </div>

    <div class="modal fade" id="banModal" tabindex="-1" aria-labelledby="banModalLabel" aria-hidden="true">
        <div class="modal-dialog modal-dialog-centered modal-dialog-scrollable">
            <div class="modal-content bg-light" id="banModalSpace">
                <div class="modal-header">
                    <h1 class="modal-title fs-5" id="banModalLabel">Ban user</h1>
                    <button type="button" class="btn-close" data-bs-dismiss="modal" aria-label="Close"></button>
                </div>
                <div class="modal-body text-dark">
                    <form id="ban-user-form" hx-post="/hx/accounts/ban_ok" hx-target="#banModalSpace" hx-trigger="submit">
                        <div class="mb-3">
                            <label for="login-input" class="col-form-label">Login:</label>
                            <input type="hidden" class="form-control" name="login" id="login-input" value="root">
                        </div>
                    </form>
                    Block the root user?
                </div>
                <div class="modal-footer">
                    <button type="button" class="btn btn-secondary" data-bs-dismiss="modal">Cancel</button>
                    <button type="submit" form="ban-user-form" class="btn btn-primary">Create</button>
                </div>
            </div>
        </div>
    </div>

    <div class="modal fade" id="unbanModal" tabindex="-1" aria-labelledby="unbanModalLabel" aria-hidden="true">
        <div class="modal-dialog modal-dialog-centered modal-dialog-scrollable">
            <div class="modal-content bg-light" id="unbanModalSpace">
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
            </div>
        </div>
    </div>

    <div class="modal fade" id="delModal" tabindex="-1" aria-labelledby="delModalLabel" aria-hidden="true">
        <div class="modal-dialog modal-dialog-centered modal-dialog-scrollable">
            <div class="modal-content bg-light" id="delModalSpace">
                <div class="modal-header">
                    <h1 class="modal-title fs-5" id="delModalLabel">Delete user</h1>
                    <button type="button" class="btn-close" data-bs-dismiss="modal" aria-label="Close"></button>
                </div>
                <div class="modal-body text-dark">
                    <form id="del-user-form" hx-post="/hx/accounts/del_ok" hx-target="#delModalSpace" hx-trigger="submit">
                        <div class="mb-3">
                            <label for="login-input" class="col-form-label">Login:</label>
                            <input type="hidden" class="form-control" name="login" id="login-input" value="root">
                        </div>
                    </form>
                    Delete the root user?
                </div>
                <div class="modal-footer">
                    <button type="button" class="btn btn-secondary" data-bs-dismiss="modal">Cancel</button>
                    <button type="submit" form="del-user-form" class="btn btn-primary">Delete</button>
                </div>
            </div>
        </div>
    </div>

    <!-- JavaScript Libraries -->
    <script src="https://code.jquery.com/jquery-3.4.1.min.js"></script>
    <script src="https://cdn.jsdelivr.net/npm/bootstrap@5.0.0/dist/js/bootstrap.bundle.min.js"></script>

    <!-- Template Javascript -->
    <script>
    `

var html4_jsblank string = `
    </script>
    <script>
    `

var html5 string = `
    </script>
</body>

</html>
`
