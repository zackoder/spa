export default class Validation {

    //Validate Nickname
    static validateNickname(nickname) {
        return this.validateName(nickname, "nickname");
    }

    //Validate Firstname
    static validateFirstname(firstname) {
        return this.validateName(firstname, "Firstname");
    }

    //Validate Lastname
    static validateLastname(lastname) {
        return this.validateName(lastname, "Lastname");
    }

    static validateName(value, name) {
        let eleErr = `err${name}`;
        let err = document.getElementById(eleErr);

        if (!value) {
            err.style.display = "block";
            err.textContent = `Please enter your ${name}.`;
            return false;
        }
        if (value.length > 20) {
            err.style.display = "block";
            err.textContent = `${name} should have 8 to 20 Letters`;
            return false;
        }
        const nicknameRegex = /^[A-Za-z][A-Za-z0-9_-]{2,20}$/;
        if (!nicknameRegex.test(value)) {
            err.style.display = "block";
            err.textContent = `Please enter your ${name} correctly`;
            return false;
        }

        err.style.display = "none";
        return true;

    }

    //Validate gender
    static validateGender(gender) {
        let errEl = document.getElementById('errgender');
        if (!gender || (gender !== "male" && gender !== "female")) {
            errEl.style.display = "block";
            errEl.textContent = "Please Check your Gender";
            return false
        }
        errEl.style.display = "none";
        return true
    }

    //Validate Age
    static validateAge(age) {
        let errEl = document.getElementById('errbirthdate');
        let dateRegex = /^\d{4}-\d{2}-\d{2}$/;

        if (!dateRegex.test(age)) {
            errEl.style.display = "block";
            errEl.textContent = "Please Check your Birth day";
            return false;
        }

        const [year, month, day] = age.split('-').map(Number);

        const date = new Date(year, month - 1, day);

        const isValidDate =
            date.getFullYear() === year &&
            date.getMonth() === month - 1 &&
            date.getDate() === day;


        const currentDate = new Date();

        let ageDiff = currentDate.getFullYear() - date.getFullYear();

        const hasBirthdayOccurred = currentDate.getMonth() > date.getMonth() ||
            (currentDate.getMonth() === date.getMonth() && currentDate.getDate() >= date.getDate());
        if (!hasBirthdayOccurred) { ageDiff-- }

        if (isValidDate && ageDiff >= 10) {
            errEl.style.display = "none";
            return true;
        } else {
            errEl.style.display = "block";
            errEl.textContent = "Please Check your Birth day";
            return false
        }

    }

    // Validate Email
    static validateEmail(email) {
        let errEl = document.getElementById('erremail');
        const emailRegex = /^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$/;
        if (!email) {
            errEl.style.display = "block";
            errEl.textContent = "Please enter your email";
            return false
        }
        if (!emailRegex.test(email)) {
            errEl.style.display = "block";
            errEl.textContent = "Please enter your email correctly";
            return false
        }
        errEl.style.display = "none";
        return true
    }

    // Validate Password
    static validatePassword(password) {
        let errEl = document.getElementById('errpassword');
        if (!password) {
            errEl.style.display = "block";
            errEl.textContent = "Please enter your password";
            return false
        }
        if (password.length < 8) {
            errEl.style.display = "block";
            errEl.textContent = "Please enter at least 8 letters";
            return false
        }

        if (!/[a-z]/.test(password)) {
            errEl.style.display = "block";
            errEl.textContent = "Please enter at least one lowercase letter";
            return false
        }
        if (!/[A-Z]/.test(password)) {
            errEl.style.display = "block";
            errEl.textContent = "Please enter at least one uppercase letter";
            return false
        }

        if (!/\d/.test(password)) {
            errEl.style.display = "block";
            errEl.textContent = "Please enter at least one digit";
            return false
        }

        if (!/[!@#$%^&*()_+\-=\[\]{};':"\\|,.<>\/?]/.test(password)) {
            errEl.style.display = "block";
            errEl.textContent = "Please enter at least one special character";
            return false
        }

        errEl.style.display = "none";
        return true
    }

    // Validate Confirm Password
    static validateConfirmPassword(password) {
        let errEl = document.getElementById('errconfirmpassword');
        let getPassword = document.getElementById('password');

        if (!password) {
            errEl.style.display = "block";
            errEl.textContent = "Please enter your password";
            return false
        }

        if (password !== getPassword.value) {
            errEl.style.display = "block";
            errEl.textContent = "your password not correct";
            return false
        }
        errEl.style.display = "none";
        return true

    }
}

