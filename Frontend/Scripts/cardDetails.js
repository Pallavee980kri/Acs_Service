var payNowButtonFlag = false;
var cardNumberValidFlag = true;
var errorMessage = "";
var cardHolderNameFlag = true;
var pleaseWaitLoading = false;
document.getElementById("cardHolderNameErrorMessage").style.display = "none";
document.getElementById("errorMessageForCvv").style.display = "none";
document.getElementById("pleaseWait").style.display = "none";
let showErrorMessage = document.getElementById("showErrorMessage");
showErrorMessage.textContent = errorMessage;
showErrorMessage.style.display = "none";
let submitButton = document.getElementById("submitButton");
let cardNumberErrorMessage = document.getElementById("cardNumberErrorMessage");
cardNumberErrorMessage.style.display = "none";
submitButton.disabled = true;

const AddExpiryYearOption = () => {
  let expiryyear = document.getElementById("expiryYear");
  for (let index = 0; index < 99; index++) {
    let option = document.createElement("option");
    option.value = 2023 + index;
    option.textContent = 2023 + index;
    expiryyear.append(option);
  }
};
AddExpiryYearOption();

const handleSubmit = async () => {
  document.getElementById("submitButton").style.display = "none";
  document.getElementById("pleaseWait").style.display = "block";
  try {
    event.preventDefault();

    let form = document.getElementById("form");
    let card_number = form.card_number.value;
    let cardHolderName = form.cardHolderName.value;
    let cvv = form.cvv.value;
    let expiryMonth = +form.expiryMonth.value;
    let expiryYear = +form.expiryYear.value;
    let formContent = {
      card_number: card_number,
      cardholder_name: cardHolderName,
      cvv: cvv,
      expiry_month: expiryMonth,
      expiry_year: expiryYear,
    };
    console.log(formContent);

    let res = await fetch("http://localhost:8000/process_payment", {
      method: "POST",
      headers: {
        "content-type": "application/json",
      },
      body: JSON.stringify(formContent),
    });
    let data = await res.json();

    if (res.status == 200) {
      document.getElementById("pleaseWait").style.display = "none";
      document.getElementById("submitButton").style.display = "block";
      localStorage.setItem("card_number", formContent.card_number);
      window.location.href = "otpPage.html";
    } else {
      // document.getElementById("pleaseWait").style.display = "block";
      document.getElementById("submitButton").style.display = "block";
      document.getElementById("pleaseWait").style.display = "none";

      errorMessage = data.error;
      showErrorMessage.textContent = errorMessage;
      showErrorMessage.style.display = "block";
    }
  } catch (err) {
    console.log(err);
  }
};
//validation for card number
const handleValidationForCardNumber = () => {
  let form = document.getElementById("form");
  let card_number = form.card_number.value;

  if (event.keyCode == 69 || event.keyCode == 187 || event.keyCode == 190) {
    event.preventDefault();
  }

  if (
    event.keyCode !== 8 &&
    event.keyCode !== 37 &&
    event.keyCode !== 39 &&
    event.keyCode !== 46
  ) {
    if (card_number.length === 16) {
      event.preventDefault();
    }
  }
};

//validation for cvv

const handleValidationForCvv = () => {
  let form = document.getElementById("form");

  let cvv = form.cvv.value;
  var key = event.key;
  var numbers = "0123456789";
  if (
    event.keyCode !== 8 &&
    event.keyCode !== 37 &&
    event.keyCode !== 39 &&
    event.keyCode !== 46
  ) {
    if (!numbers.includes(key)) {
      event.preventDefault();
    }
  }
  if (
    event.keyCode !== 8 &&
    event.keyCode !== 37 &&
    event.keyCode !== 39 &&
    event.keyCode !== 46
  ) {
    if (cvv.length === 3) {
      event.preventDefault();
    }
  }
};

const hanleCheckPayNowbuttonEnable = (type) => {
  showErrorMessage.style.display = "none";

  let form = document.getElementById("form");
  let errorMessageBox = document.getElementById("errorMessageForCvv");
  let card_number = form.card_number.value;
  let cardHolderName = form.cardHolderName.value;
  let cvv = form.cvv.value;
  if (type == "cardNumberInput") {
    if (card_number.length == 16) {
      let cardNumberErrorMessage = document.getElementById(
        "cardNumberErrorMessage"
      );
      cardNumberErrorMessage.style.display = "none";
      cardNumberValidFlag = true;
    } else {
      cardNumberValidFlag = false;
      let cardNumberErrorMessage = document.getElementById(
        "cardNumberErrorMessage"
      );
      cardNumberErrorMessage.style.display = "block";
    }
  } else if (type == "carhHolderNameInput") {
    if (cardHolderName.trim().length == 0) {
      cardHolderNameFlag = false;
      document.getElementById("cardHolderNameErrorMessage").style.display =
        "block";
    } else {
      cardHolderNameFlag = true;
      document.getElementById("cardHolderNameErrorMessage").style.display =
        "none";
    }
  } else if (type == "cvvInput") {
    if (cvv.length < 3) {
      errorMessageBox.style.display = "block";
    } else {
      errorMessageBox.style.display = "none";
    }
  }

  if (
    card_number.length == 16 &&
    cardHolderName.trim().length != 0 &&
    cvv.length == 3
  ) {
    payNowButtonFlag = true;
  } else {
    payNowButtonFlag = false;
  }

  if (payNowButtonFlag) {
    let submitButton = document.getElementById("submitButton");
    submitButton.style.backgroundColor = "purple";
    submitButton.style.opacity = "1";
    submitButton.disabled = false;
  } else {
    let submitButton = document.getElementById("submitButton");
    submitButton.style.backgroundColor = "purple";
    submitButton.style.opacity = "0.5";
    submitButton.disabled = true;
  }
};
