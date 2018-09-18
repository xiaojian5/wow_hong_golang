var vue = new Vue({
    el: "#fast-create",
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
        skillTypeTable: [
            {
                "id": 0,
                "name": "请选择技能类型",
            },
            {
                "id": 1,
                "name": "伤害技能(单体)",
            },
            {
                "id": 2,
                "name": "治疗技能(单体)",
            },
            {
                "id": 3,
                "name": "范围(伤害/治疗)技能",
            },
            {
                "id": 4,
                "name": "增益技能",
            },

        ],
        skillName: "",
        skillType: "0",
        skillCondition: {
            "0": [],
            "1": [
                [
                    {
                        "condition": "[@mouseover,harm]",
                        "desc": "鼠标指向的目标是敌对目标",
                    },
                    {
                        "condition": "[@target,harm]",
                        "desc": "当前目标是敌对目标",
                    },
                    {
                        "condition": "",
                        "desc": "直接施放"
                    },
                ],
            ],
            "2": [
                [
                    {
                        "condition": "[@mouseover,help]",
                        "desc": "鼠标指向的目标是友善目标",
                    },
                    {
                        "condition": "[@target,help]",
                        "desc": "当前目标是友善目标",
                    },
                    {
                        "condition": "[@player]",
                        "desc": "以玩家自己为目标",
                    },
                    {
                        "condition": "",
                        "desc": "直接施放"
                    }
                ]
            ],
            "3": [
                [
                    {
                        "condition": "[@cursor]",
                        "desc": "鼠标当前位置",
                    },
                    {
                        "condition": "",
                        "desc": "直接施放"
                    }
                ],
                [
                    {
                        "condition": "[@player]",
                        "desc": "玩家当前位置",
                    },
                    {
                        "condition": "",
                        "desc": "直接施放"
                    }
                ]
            ],
            "4": [
                [
                    {
                        "condition": "",
                        "desc": "直接施放"
                    }
                ],
            ]
        },
        skillTable: [],
        skillTableOfRemote: [],
    },
    methods: {
        fastCreate: function() {
            console.log(vue.$data.skillType);
            if (vue.$data.skillType !== "0" && vue.$data !== 0) {
                vue.$data.skillName = vue.$data.skillName.replace(/(^\s*)|(\s*$)/g, "");
                vue.$data.skillTable = [];//清空
                vue.$data.skillTableOfRemote = [];
                var num = 1;

                var skills = vue.$data.skillCondition[vue.$data.skillType];

                var len = skills.length;
                for (var i = 0; i < len; i++) {
                    var text = desc = "";
                    var lenL = skills[i].length;
                    for (var j = 0; j < lenL; j++) {
                        text += "<br>/cast " + skills[i][j].condition + vue.$data.skillName;
                        if (skills[i][j].condition !== "") {

                            desc += "<br>- 【" + skills[i][j].desc + "】则施放【" + vue.$data.skillName +
                                "】，否则跳过";
                        } else {
                            desc += "<br>- " + skills[i][j].desc + "【" + vue.$data.skillName + "】";

                        }
                    }
                    vue.$data.skillTable.push({
                        "id": num,
                        "text": "#showtooltip" + text,
                        "desc": "- 【显示技能图标和技能描述】" + desc,
                    });
                    num++;
                }
                axios.get('/macros', {
                    params: {
                        macro: vue.$data.skillName,
                    }
                })
                    .then(function(response) {
                        // console.log(response);
                        var len = response.data.length
                        for (var i = 0; i < len; i++) {
                            vue.$data.skillTableOfRemote.push({
                                "id": num,
                                "text": response.data[i].macro,
                                "desc": response.data[i].title,
                            });
                            num++;
                        }
                    })
                    .catch(function(error) {
                        console.log(error);
                    });
            } else {
                alert("请选择技能类型")
            }
        }
    },
    computed: {
    },
    created: function() {
        axios.post('/log/index', {
        })
            .then(function(response) {
                // console.log(response);
            })
            .catch(function(error) {
                console.log(error);
            });
    }
});