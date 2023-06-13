//validation for the OTP
function handleValidationForOtp() {
  let otp = event.target.value;
  if (event.keyCode == 69 || event.keyCode == 187 || event.keyCode == 190) {
    event.preventDefault();
  }
  if (event.keyCode !== 8) {
    if (otp.length === 6) {
      event.preventDefault();
    }
  }
}
handleValidationForOtp();
