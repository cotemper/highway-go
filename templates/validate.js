function validateForm() {
    // Grab registerInput value
    let x = document.forms['registerForm']['registerInput'].value;
    let special_char_regex = /[^A-Za-z 0-9]/g;

    // Character length guard
    if (x.length < 2 || x.length > 20) {
        alert('Must be between 2-20 characters.');
        return false;
    }

    // Special character guard
    if (special_char_regex.test(x)) {
        alert('Special characters are not allowed.');
        return false;
    }

    makeCredential()

    return true;
}