function change() {
  var email = document.getElementById("Email").value;
  if (email != "") {
    document.getElementById("continue").style.display = "none";
    document.getElementById("submitForm2").style.display = "block";
    document.getElementById("inputForm").style.display = "none";
    document.getElementById("inputForm2").style.display = "block";
    document.getElementById("spanBadLogin").style.display = "none";
    document.getElementById("confirmation_password").style.display = "flex";
    document.getElementById("emailLogin").innerHTML = email;
    document.getElementById("lbltipAddedComment").style.display = "none";
  } else {
    document.getElementById("lbltipAddedComment").innerHTML =
      "The Email field cannot be empty";
    document.getElementById("lbltipAddedComment").style.display = "block";
  }
}
function back() {
  document.getElementById("continue").style.display = "block";
  document.getElementById("submitForm").style.display = "none";
  document.getElementById("inputForm").style.display = "block";
  document.getElementById("inputForm2").style.display = "none";
  document.getElementById("spanBadLogin").style.display = "none";
  document.getElementById("lbltipAddedComment").style.display = "none";
  document.getElementById("confirmation_password").style.display = "none";
}
