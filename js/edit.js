var vue = new Vue({
    el: "#edit",
    data: {
        professions: [
            {"id": 1, "name": "猎人"},
            {"id": 10, "name": "德鲁伊"},
            {"id": 20, "name": "死亡骑士"},
            {"id": 30, "name": "法师"},
            {"id": 40, "name": "圣骑士"},
            {"id": 50, "name": "牧师"},
            {"id": 60, "name": "盗贼"},
            {"id": 70, "name": "萨满祭司"},
            {"id": 80, "name": "术士"},
            {"id": 90, "name": "战士"},
            {"id": 100, "name": "武僧"},
            {"id": 110, "name": "恶魔猎手"},
            {"id": 120, "name": "其他"},
        ],
        macros: [],
    },
    methods: {
        search: function() {
            var verify = document.getElementById("verify")
            isVerify = verify.value;
            axios.get('/macros', {
                params: {
                    isVerify: isVerify,
                }
            })
                .then(function(response) {
                    console.log(response);
                    vue.$data.macros = response.data;
                })
                .catch(function(error) {
                    console.log(error);
                })
                .then(function() {
                    // always executed
                });
        },
        update: function(id) {
            console.log(id);
            var isVerify = document.getElementById("verify_" + id).value
            var author = document.getElementById("author_" + id).value
            var title = document.getElementById("title_" + id).value
            var macro = document.getElementById("macro_" + id).value

            axios.put('/macros', {
                id: parseInt(id),
                isVerify: parseInt(isVerify),
                author: author,
                title: title,
                macro: macro,
            })
                .then(function(response) {
                    console.log(response);
                })
                .catch(function(error) {
                    console.log(error);
                })
                .then(function() {
                    // always executed
                });
        }
    },
})