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
                {
                    "id": 1,
                    "condition": "[@mouseover,harm]",
                    "desc": "鼠标指向的目标是敌对目标",
                },
                {
                    "id": 2,
                    "condition": "[@target,harm]",
                    "desc": "当前目标是敌对目标",
                },
            ],
            "2": [
                {
                    "id": 3,
                    "condition": "[@mouseover,help]",
                    "desc": "鼠标指向的目标是友善目标",
                },
                {
                    "id": 4,
                    "condition": "[@target,help]",
                    "desc": "当前目标是友善目标",
                },
                {
                    "id": 5,
                    "condition": "[@player]",
                    "desc": "以玩家自己为目标",
                },
            ],
            "3": [
                {
                    "id": 6,
                    "condition": "[@player]",
                    "desc": "玩家当前位置",
                },
                {
                    "id": 7,
                    "condition": "[@cursor]",
                    "desc": "鼠标当前位置",
                },
            ],
        },
        skillTable: [],
    },
    methods: {
        fastCreate: function() {
            console.log(vue.$data.skillType);
            if (self.skillType !== "0") {
                vue.$data.skillTable = [];//清空

                var skills = vue.$data.skillCondition[vue.$data.skillType];

                var text = desc = "";
                var len = skills.length;
                for (var i = 0; i < len; i++) {
                    text += "<br>/cast " + skills[i].condition + vue.$data.skillName;
                    desc += "<br>- 【" + skills[i].desc + "】则施放【" + vue.$data.skillName + "】，否则跳过";
                }
                vue.$data.skillTable.push({
                    "id": 1,
                    "text": "#showtooltip" + text,
                    "desc": "- 【显示技能图标和技能描述】" + desc,
                });
            }
        }
    },
    computed: {},
    created: function() {
        console.log("create");
    }
});