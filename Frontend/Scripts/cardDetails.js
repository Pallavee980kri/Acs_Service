//creating options for Expiry Year
var payNowButtonFlag = false;
var cardNumberValidFlag = true;

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

const handleSubmit = () => {
  event.preventDefault();

  let form = document.getElementById("form");
  let cardNumber = form.cardNumber.value;
  let cardHolderName = form.cardHolderName.value;
  let cvv = form.cvv.value;
  let expiryMonth = form.expiryMonth.value;
  let expiryYear = form.expiryYear.value;
  let formContent = {
    cardNumber,
    cardHolderName,
    cvv,
    expiryMonth,
    expiryYear,
  };
  console.log(formContent);
};
//validation for card number
const handleValidationForCardNumber = () => {
  let form = document.getElementById("form");
  let cardNumber = form.cardNumber.value;

  if (
    event.keyCode == 69 ||
    event.keyCode == 187 ||
    event.keyCode == 190 ||
    event.keyCode == 189
  ) {
    event.preventDefault();
  }

  if (
    event.keyCode !== 8 &&
    event.keyCode !== 37 &&
    event.keyCode !== 39 &&
    event.keyCode !== 46
  ) {
    if (cardNumber.length === 16) {
      event.preventDefault();
    }
  }
  if (cardNumber.length == 16) {
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

const hanleCheckPayNowbuttonEnable = () => {
  let form = document.getElementById("form");
  let cardNumber = form.cardNumber.value;
  let cardHolderName = form.cardHolderName.value;
  let cvv = form.cvv.value;
  if (
    cardNumber.length == 16 &&
    cardHolderName.length != 0 &&
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
