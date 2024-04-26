// Background light
const cursor = document.querySelector('#cursor');
const onMouseChangePosition = (event) => {
    cursor.style.left = event.clientX + 'px';
    cursor.style.top = event.clientY + 'px';
};
onmousemove = onMouseChangePosition;

// Logout modal
htmx.on("#uBtnLogOut", "click", function(){ 
    htmx.removeClass(htmx.find("#logoutModal"), "hidden");
    htmx.addClass(htmx.find('#logoutModal'), 'showed');
});
htmx.on("#logoutModalBtnClose", "click", function(){ 
    htmx.removeClass(htmx.find("#logoutModal"), "showed");
    htmx.addClass(htmx.find('#logoutModal'), 'hidden');
});
htmx.on("#logoutModalCancel", "click", function(){ 
    htmx.removeClass(htmx.find("#logoutModal"), "showed");
    htmx.addClass(htmx.find('#logoutModal'), 'hidden');
});

// Selfedit (Profile) modal
function showProfileModal() {
    htmx.removeClass(htmx.find("#profileModal"), "hidden");
    htmx.addClass(htmx.find("#profileModal"), "showed");
}
function hideProfileModal() {
    htmx.removeClass(htmx.find("#profileModal"), "showed");
    htmx.addClass(htmx.find("#profileModal"), 'hidden');
}

// Create modal
function showCreateModal() {
    htmx.removeClass(htmx.find("#createModal"), "hidden");
    htmx.addClass(htmx.find("#createModal"), "showed");
}
function hideCreateModal() {
    htmx.removeClass(htmx.find("#createModal"), "showed");
    htmx.addClass(htmx.find("#createModal"), 'hidden');
}

// Edit modal
function showEditModal() {
    htmx.removeClass(htmx.find("#editModal"), "hidden");
    htmx.addClass(htmx.find("#editModal"), "showed");
}
function hideEditModal() {
    htmx.removeClass(htmx.find("#editModal"), "showed");
    htmx.addClass(htmx.find("#editModal"), 'hidden');
}

// Ban modal
function showBanModal() {
    htmx.removeClass(htmx.find("#banModal"), "hidden");
    htmx.addClass(htmx.find("#banModal"), "showed");
}
function hideBanModal() {
    htmx.removeClass(htmx.find("#banModal"), "showed");
    htmx.addClass(htmx.find("#banModal"), 'hidden');
}

// UnBan modal
function showUnBanModal() {
    htmx.removeClass(htmx.find("#unbanModal"), "hidden");
    htmx.addClass(htmx.find("#unbanModal"), "showed");
}
function hideUnBanModal() {
    htmx.removeClass(htmx.find("#unbanModal"), "showed");
    htmx.addClass(htmx.find("#unbanModal"), 'hidden');
}

// Delete modal
function showDelModal() {
    htmx.removeClass(htmx.find("#delModal"), "hidden");
    htmx.addClass(htmx.find("#delModal"), "showed");
}
function hideDelModal() {
    htmx.removeClass(htmx.find("#delModal"), "showed");
    htmx.addClass(htmx.find("#delModal"), 'hidden');
}
