let vue = new Vue({
	el: "#container",
	delimiters: ['${', '}$'],
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
			// console.log('insertMacro_'+vue.$datacondition_radio);
			vue.$data.skillName = vue.$data.skillName.replace(/(^\s*)|(\s*$)/g, "");
			//自己组合条件
			if (vue.$data.condition_radio === 2) {
				//分割多选条件
				vue.$data.selectConditions = '[' + vue.$data.select_item + ']';

				let changeLine = '';//是否换行
				//如果是castsequence是特殊的
				if (vue.$data.selectCommand === '/castsequence') {
					let s = ' reset=12/combat/target';
					if (vue.$data.macro) {
						changeLine = "\r\n";
					}
					vue.$data.macro += changeLine + vue.$data.selectCommand + s + ' ' +
						vue.$data.skillName;
				} else {
					if (vue.$data.macro) {
						let ok = confirm('是否换行');
						if (ok) {
							changeLine = "\r\n";
							vue.$data.macro += changeLine + vue.$data.selectCommand + ' ' +
								vue.$data.selectConditions + vue.$data.skillName;
						} else {
							changeLine = ';';
							vue.$data.macro += changeLine + vue.$data.selectConditions +
								vue.$data.skillName;//不换行不加命令
						}
					} else {
						if (!vue.$data.selectCommand) {
							alert('请选择命令');
						} else {
							vue.$data.macro = changeLine + vue.$data.selectCommand + ' ' +
								vue.$data.selectConditions + vue.$data.skillName;
						}
					}
				}
			} else {
				vue.$data.macro = '';
				vue.$data.selectTemplates = vue.$data.select_template.split(';');
				// console.log(vue.$data.selectTemplates);
				let len = vue.$data.selectTemplates.length;
				let splitStr = '';
				for (let i = 0; i < len; i++) {
					if (vue.$data.macro) {
						splitStr = ';';
					}
					// console.log(vue.$data.selectTemplates[i]);
					vue.$data.macro += splitStr + vue.$data.selectTemplates[i] + vue.$data.skillName;
				}
			}
			vue.$data.allMacro = "#showtooltip\r\n" + vue.$data.macro;
		},
		copyHong: function() {
			document.getElementById("hong_text").select();
			document.execCommand("Copy");
			alert("已复制，可以进入游戏在宏命令中新建宏，然后CTRL+V粘贴进去了哦^-^!");
		},
		changeCommand: function() {
			vue.$data.selectCommand = vue.$data.commands[vue.$data.selectCommandKey].name;
			let key = vue.$data.selectCommandKey;
			let val = vue.$data.commands[key];

			let con = document.getElementById("condition");
			if (val.type !== '1') {
				con.style.display = "none";
			} else {
				con.style.display = "block";
			}
		}
	},
	computed: {},
	created: function() {
		axios.post('/log/macroIndex', {})
			.then(function(response) {
				// console.log(response);
			});
	}
});