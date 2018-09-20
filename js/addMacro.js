var vue = new Vue({
    el: "#container",
    data: {
        professions: [
            {'id': 1, 'name': '猎人'},
            {'id': 10, 'name': '德鲁伊'},
            {'id': 20, 'name': '死亡骑士'},
            {'id': 30, 'name': '法师'},
            {'id': 40, 'name': '圣骑士'},
            {'id': 50, 'name': '牧师'},
            {'id': 60, 'name': '盗贼'},
            {'id': 70, 'name': '萨满祭司'},
            {'id': 80, 'name': '术士'},
            {'id': 90, 'name': '战士'},
            {'id': 100, 'name': '武僧'},
            {'id': 110, 'name': '恶魔猎手'},
            {'id': 120, 'name': '其他'},
        ],
        select_proid: '',
        select_pro: [],
        select_mastery: [],
    },
    methods: {
        changeProfession: function() {
            var len = vue.$data.professions.length;
            console.log(len);
            if (len <= 0) {
                return false;
            }
            for (var i = 0; i < len; i++) {
                if (vue.$data.select_proid === vue.$data.professions[i].id) {
                    vue.$data.select_pro = vue.$data.professions[i];
                    // this.select_mastery = this.select_pro.child;
                }
            }
        },
        submitToService: function() {
            var pid = document.getElementById("profession").value;
            if (!pid) {
                alert('职业不能为空！');
                return false;
            }
            var title = document.getElementById("title").value;
            var macro = document.getElementById("macro").value;
            if (!title || !macro) {
                alert('标题或宏模板不能为空！');
                return false;
            }

            var author = document.getElementById("author").value;
            if (!author) {
                alert('请留下你的艾泽拉斯名称，分享者必将被铭记！');
                return false;
            }
            axios.post('/macros', {
                professionId: parseInt(pid),
                title: title,
                author: author,
                macro: macro,
            })
                .then(function(response) {
                    console.log(response);
                    document.getElementById("title").value = "";
                    document.getElementById("macro").value = "";
                })
                .catch(function(error) {
                    console.log(error);
                })
                .then(function() {
                    // always executed
                });
        }
    },
    computed: {
        select_mastery: function() {
            if (vue.$data.select_pro) {
                return vue.$data.select_pro.child;
            }
            return [];
        }
    },
    created: function() {
        axios.post('/log/addMacro', {
        })
            .then(function(response) {
                console.log(response);
            })
            .catch(function(error) {
                console.log(error);
            });
    }
});
