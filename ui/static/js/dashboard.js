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

jQuery(document).ready(function(){
    // $("#uBtnLogOut").click(function() {
    //     $("#logoutModal").removeClass("hidden").addClass("showed");
    // });
    // $("#logoutModalBtnClose").click(function() {
    //     $("#logoutModal").removeClass("showed").addClass("hidden");
    // });
    // $("#logoutModalCancel").click(function() {
    //     $("#logoutModal").removeClass("showed").addClass("hidden");
    // });

    $("#uBtnYourProfile").click(function() {
        $("#profileModal").removeClass("hidden").addClass("showed");
    });
    $("#profileModalBtnClose").click(function() {
        $("#profileModal").removeClass("showed").addClass("hidden");
    });
    $("#profileModalCancel").click(function() {
        $("#profileModal").removeClass("showed").addClass("hidden");
    });

    // $("#aBtnCreateAccount").click(function() {
    //     $("#createModal").removeClass("hidden").addClass("showed");
    // });
    // $("#createModalBtnClose").click(function() {
    //     $("#createModal").removeClass("showed").addClass("hidden");
    // });
    // $("#createModalCancel").click(function() {
    //     $("#createModal").removeClass("showed").addClass("hidden");
    // });

    // $("#editModalBtnClose").click(function() {
    //     $("#editModal").removeClass("showed").addClass("hidden");
    // });
    // $("#editModalCancel").click(function() {
    //     $("#editModal").removeClass("showed").addClass("hidden");
    // });

    // $("#banModalBtnClose").click(function() {
    //     $("#banModal").removeClass("showed").addClass("hidden");
    // });
    // $("#banModalCancel").click(function() {
    //     $("#banModal").removeClass("showed").addClass("hidden");
    // });
    
    // $("#unbanModalBtnClose").click(function() {
    //     $("#unbanModal").removeClass("showed").addClass("hidden");
    // });
    // $("#unbanModalCancel").click(function() {
    //     $("#unbanModal").removeClass("showed").addClass("hidden");
    // });

    $("#delModalBtnClose").click(function() {
        $("#delModal").removeClass("showed").addClass("hidden");
    });
    $("#delModalCancel").click(function() {
        $("#delModal").removeClass("showed").addClass("hidden");
    });

    // $(".trigger-edit").click(function() {
    //     $("#editModal").removeClass("hidden").addClass("showed");
    // });
    // $(".trigger-ban").click(function() {
    //     $("#banModal").removeClass("hidden").addClass("showed");
    // });
    // $(".trigger-unban").click(function() {
    //     $("#unbanModal").removeClass("hidden").addClass("showed");
    // });
    // $(".trigger-del").click(function() {
    //     $("#delModal").removeClass("hidden").addClass("showed");
    // });

  });