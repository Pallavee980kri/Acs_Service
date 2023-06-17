let errorMessage = "";
let showErrorMessage = document.getElementById("showErrorMessage");
showErrorMessage.textContent = errorMessage;
document.getElementById("pleaseWait").style.display = "none";
document.getElementById("pleaseWaitForResendOTP").style.display="none"
let cardNumber = localStorage.getItem("card_number");

const handleSubmitOtp = async () => {
  document.getElementById("sendOtp").style.display = "none";
  document.getElementById("pleaseWait").style.display = "block";
  // document.getElementById("pleaseWaitForResendOTP").style.display="block"
  const payload = {
    card_number: cardNumber,
    OTP: +document.getElementById("otp").value,
  };
  try {
    let res = await fetch("http://localhost:8000/match_otp", {
      method: "POST",
      headers: {
        "content-type": "application/json",
      },
      body: JSON.stringify(payload),
    });
    let data = await res.json();
    // console.log(res)
    // console.log(data)
    if (res.status == 200) {
      document.getElementById("pleaseWait").style.display = "none";
      document.getElementById("sendOtp").style.display = "block";
      alert("Payment Successfully Done !");
    } else {
      document.getElementById("sendOtp").style.display = "block";
      errorMessage = data.error;
      showErrorMessage.textContent = errorMessage;
    }
  } catch (err) {
    console.log(err);
  }
};

//validation for the OTP
function handleValidationForOtp() {
  let otp = event.target.value;
  console.log(otp);
  if (event.keyCode == 69 || event.keyCode == 187 || event.keyCode == 190) {
    event.preventDefault();
  }
  if (event.keyCode !== 8) {
    if (otp.length === 6) {
      event.preventDefault();
    }
  }
}
// handleValidationForOtp();

const handleResendOtp=async()=>{
  document.getElementById("pleaseWaitForResendOTP").style.display="block"
  document.getElementById("resendOtp").style.display="none"
  const payload={
    card_number:cardNumber
  }
  try{
let res = await fetch("http://localhost:8000/resend_otp",{
  method:"POST",
  headers:{
    "content-type":"application/json"
  },
  body:JSON.stringify(payload)
})
let data=await res.json();
if(res.status==200){
  document.getElementById("pleaseWaitForResendOTP").style.display="none"
  document.getElementById("resendOtp").style.display="block"
  alert("OTP resend Successfully...")
}
else{
  document.getElementById("resendOtp").style.display="block"
  document.getElementById("pleaseWaitForResendOTP").style.display="none"
  errorMessage = data.error;
  showErrorMessage.textContent = errorMessage;
}
  }
  catch(err){
    console.log(err)
  }

}