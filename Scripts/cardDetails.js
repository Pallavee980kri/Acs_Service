//creating options for Expiry Year

function AddExpiryYearOption() {
  let expiryyear = document.getElementById("expiryYear");
  let arr = new Array(100).fill(1);
  arr.map((element, index) => {
    let option = document.createElement("option");
    option.value = 2023 + index;
    option.textContent = 2023 + index;
    expiryyear.append(option);
  });
}
AddExpiryYearOption();

function handleSubmit() {
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
}

function handleValidationForCardNumber() {
  let form = document.getElementById("form");
  let cardNumber = form.cardNumber.value;

  if (event.keyCode == 69 || event.keyCode == 187 || event.keyCode == 190) {
    event.preventDefault();
  }

  if (event.keyCode !== 8) {
    if (cardNumber.length === 16) {
      event.preventDefault();
    }
  }
}
