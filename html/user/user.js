function loadUsers(){
    $(function(){
        $("#userDg").iDatagrid({
            idField: 'UserID',
            sortOrder:'desc',
            sortName:'CreateTime',
            pageSize:50,
            frozenColumns:[[
                {field: 'ck', title: '', checkbox: true},
                {field: 'op', title: '操作', sortable: false, halign:'center',align:'left',
                    width1: 120, formatter:userOperateFormatter},
                {field: 'UserName', title: '用户名', sortable: true,
                    formatter:$.iGrid.tooltipformatter(),
                    width: 180},
                {field: 'UserID', title: '用户ID', sortable: true,
                    formatter:$.iGrid.tooltipformatter(),
                    width: 260},
            ]],
            onBeforeLoad:function (param){
                refreshUsers(param)
            },
            columns: [[
                {field: 'CreateTime', title: '创建时间', sortable: true,
                    formatter:$.iGrid.tooltipformatter(),width: 260}
            ]]
        });
    });
}

function userOperateFormatter(value, row, index) {
    let htmlstr = "";
    htmlstr += '<button class="layui-btn-blue layui-btn layui-btn-xs" onclick="resetUserPwd(\'' + row.UserID + '\')">改密</button>';
    htmlstr += '<button class="layui-btn-gray layui-btn layui-btn-xs" onclick="removeUser(\'' + row.UserID + '\')">删除</button>';
    return htmlstr;
}

function refreshUsers(param){
    $.app.getJson(V3_API_URL + '/user/list', null, function (data) {
        if(data.status !== 0){
            $.app.alert(data.msg || '获取用户列表失败');
            return;
        }

        let rows = data.data || [];
        $('#userDg').datagrid('loadData', {
            total: rows.length,
            rows: rows
        })

    }, true);
}

function createUser(){
    $.iDialog.openDialog({
        title: '新增用户',
        minimizable:false,
        id:'addUserDlg',
        width: 560,
        height: 300,
        href:'./add_user.html',
        buttonsGroup: [{
            text: '确定',
            iconCls: 'fa fa-floppy-o',
            btnCls: 'cubeui-btn-orange',
            handler:'ajaxForm',
            beforeAjax:function(o){
                let t = this;
				// o.ajaxData is querystring, convert to json
				let params = $.extends.json.param2json(o.ajaxData);

                $.app.post(V3_API_URL + '/user/create', params, function(resp){
                    if(resp.status === 0){
                        $.app.show('创建用户成功');
                        $.iDialog.closeOutterDialog($(t))
                        reloadUserDg();
                    }else{
                        $.app.alert(resp.msg || '创建用户失败');
                    }
                })

                return false;
            }
        }]
    });
}

function removeUser(userid){
    if(userid==null){
        let rows = $('#userDg').datagrid('getChecked');

        if(rows.length>1){
            $.app.show('本版本仅支持选择一个用户删除');
            return ;
        }

        if(rows.length==0){
            $.app.show('请选择一个用户删除');
            return;
        }else{
            userid = rows[0].UserID;
        }
    }

    $.app.confirm("删除用户", "您确定要删除所选择的用户？",function () {
        $.app.post(V3_API_URL + '/user/delete', {userid: userid}, function (resp) {
            if(resp.status === 0){
                $.app.show("删除用户成功");
                reloadUserDg();
            }else{
                $.app.alert(resp.msg || '删除用户失败');
            }
        })
    })
}

function resetUserPwd(userid){
    if(userid==null){
        let rows = $('#userDg').datagrid('getChecked');

        if(rows.length>1){
            $.app.show('本版本仅支持选择一个用户改密');
            return ;
        }

        if(rows.length==0){
            $.app.show('请选择一个用户改密');
            return;
        }else{
            userid = rows[0].UserID;
        }
    }

    $.iDialog.openDialog({
        title: '修改用户密码',
        minimizable:false,
        id:'resetPwdDlg',
        width: 560,
        height: 240,
        href:'./resetpwd.html',
        buttonsGroup: [{
            text: '确定',
            iconCls: 'fa fa-floppy-o',
            btnCls: 'cubeui-btn-orange',
            handler:'ajaxForm',
            beforeAjax:function(o){
                let t = this;
				let params = $.extends.json.param2json(o.ajaxData);
				params.userid = userid;

                $.app.post(V3_API_URL + '/user/resetpwd', params, function(resp){
                    if(resp.status === 0){
                        $.app.show('修改密码成功');
                        $.iDialog.closeOutterDialog($(t))
                        reloadUserDg();
                    }else{
                        $.app.alert(resp.msg || '修改密码失败');
                    }
                })

                return false;
            }
        }]
    });
}

function reloadUserDg(){
    $('#userDg').datagrid('reload');
    $('#layout').layout('resize');
}

function onActivated(opts, title, idx){
    reloadUserDg();
}
