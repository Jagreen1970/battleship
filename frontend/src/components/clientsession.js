
const UserSession = (function () {
    let userName = "";

    let getUserName = function () {
        return userName;
    }

    let setUserName = function (name) {
        userName = name;
    }

    return {
        getUserName: getUserName,
        setUserName: setUserName
    }

})();

export default UserSession;