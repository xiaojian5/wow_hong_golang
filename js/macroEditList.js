var vue = new Vue({
    el: "#macrolist",
    data: {
        professions: [
            {'id':1, 'name':'猎人'},
            {'id':10, 'name':'德鲁伊'},
            {'id':20, 'name':'死亡骑士'},
            {'id':30, 'name':'法师'},
            {'id':40, 'name':'圣骑士'},
            {'id':50, 'name':'牧师'},
            {'id':60, 'name':'盗贼'},
            {'id':70, 'name':'萨满祭司'},
            {'id':80, 'name':'术士'},
            {'id':90, 'name':'战士'},
            {'id':100, 'name':'武僧'},
            {'id':110, 'name':'恶魔猎手'},
            {'id':120, 'name':'其他'},
        ],
        macros: [],
    },
    methods: {
        search: function() {
            var pid = $("#profession").val();
            var verify = $("#verify").val();

            if(!pid){
                alert('职业或专精未选择！');
                return false;
            }

            $.getJSON('./action.php?method=getMacros',{pid:pid,verify:verify},function (json) {
                this.macros = json;
                // this.$set('macros', json);
            }.bind(this));
        },
        update: function(event) {
            var id = event.target.attributes.flag.value;

            var author = $("#author_"+id).val();
            var title = $("#title_"+id).val();
            var macro = $("#macro_"+id).val();
            var verify = $("#verify_"+id).val();

            $.post('./action.php?method=updateMacro',{id:id,author:author,title:title,macro:macro,verify:verify},function (result) {
                console.log(result);
                if(result){
                    alert("修改成功！");
                }else{
                    alert('修改失败！');
                }
            }.bind(this));
        },
        deleteMacro: function(item, event) {
            var id = event.target.attributes.flag.value;

            var ok = confirm('是否删除？');
            if(ok){
                $.getJSON('./action.php?method=deleteMacro',{id:id},function (result) {
                    if(result){
                        // this.macros.$remove(item);
                        var index = this.macros.indexOf(item);
                        this.macros.splice(index, 1);
                    }else{
                        alert('删除失败！');
                    }
                }.bind(this));
            }
        }
    },
    created: function(){
    },
    ready: function() {
    }
});
