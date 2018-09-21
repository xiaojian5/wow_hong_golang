let vue = new Vue({
    el: "#container",
    data: {
        commands: [],
        conditions: [],
        templates: [],
        condition_radio: 2,
        selectCommandKey: '',
        selectCommand: '',
        select_item: '',//多选条件，需要分割开
        selectConditions: [],
        select_template: '',//组合模板是很多条件
        selectTemplates: [],
        skillName: '',
        macro: '',//组合条件组合成的宏
        allMacro: ''//加上头部tip的宏
    },
    methods: {
        insertMacro: function() {
            // console.log('insertMacro_'+this.condition_radio);
            this.skillName = this.skillName.replace(/(^\s*)|(\s*$)/g, "");
            //自己组合条件
            if (this.condition_radio === 2) {
                //分割多选条件
                this.selectConditions = '[' + this.select_item + ']';

                let changeLine = '';//是否换行
                //如果是castsequence是特殊的
                if (this.selectCommand === '/castsequence') {
                    let s = ' reset=12/combat/target';
                    if (this.macro) {
                        changeLine = "\r\n";
                    }
                    this.macro += changeLine + this.selectCommand + s + ' ' + this.skillName;
                } else {
                    if (this.macro) {
                        let ok = confirm('是否换行');
                        if (ok) {
                            changeLine = "\r\n";
                            this.macro += changeLine + this.selectCommand + ' ' +
                                this.selectConditions + this.skillName;
                        } else {
                            changeLine = ';';
                            this.macro += changeLine + this.selectConditions + this.skillName;//不换行不加命令
                        }
                    } else {
                        if (!this.selectCommand) {
                            alert('请选择命令');
                        } else {
                            this.macro = changeLine + this.selectCommand + ' ' +
                                this.selectConditions +
                                this.skillName;
                        }
                    }
                }
            } else {
                this.macro = '';
                this.selectTemplates = this.select_template.split(';');
                // console.log(this.selectTemplates);
                let len = this.selectTemplates.length;
                let splitStr = '';
                for (let i = 0; i < len; i++) {
                    if (this.macro) {
                        splitStr = ';';
                    }
                    // console.log(this.selectTemplates[i]);
                    this.macro += splitStr + this.selectTemplates[i] + this.skillName;
                }
            }
            this.allMacro = "#showtooltip\r\n" + this.macro;
        },
        copyHong: function() {
            document.getElementById("hong_text").select();
            document.execCommand("Copy");
            alert("已复制，可以进入游戏在宏命令中新建宏，然后CTRL+V粘贴进去了哦^-^!");
        },
        changeCommand: function(event) {
            this.selectCommand = this.commands[this.selectCommandKey].name;
            let key = this.selectCommandKey;
            let val = this.commands[key];

            let con = document.getElementById("condition");
            if (val.type === '1') {
                con.style.display = "none";
            } else {
                con.style.display = "block";
            }
        }
    },
    computed: {},
    created: function() {
        axios.post('/log/macroIndex', {
        })
            .then(function(response) {
                // console.log(response);
            });
    }
});