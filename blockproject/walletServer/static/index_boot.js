$("#reload_wallet").click(function () {
    $.ajax({
        url: "http://127.0.0.1:8080/wallet",
        type: "POST",
        success: function (response) {
            $("#inputPublic").val(response["public_key"]);
            $("#inputPrivateKey").val(response["private_key"]);
            $("#inputAddress").val(response["blockchain_address"]);
            console.info(response);
        },
        error: function (error) {
            console.error(error);
        },
    })
})

$("#reload").click(function () {
    var privateKeyValue = $("#inputPrivateKey").val();
    $.ajax({
        url: "http://127.0.0.1:8080/walletByPrivatekey",
        type: "POST",
        data: {
            privatekey: privateKeyValue,
        },
        success: function (response) {
            $("#inputPublic").val(response["public_key"]);
            $("#inputPrivateKey").val(response["private_key"]);
            $("#inputAddress").val(response["blockchain_address"]);
            console.info(response);
        },
        error: function (error) {
            console.error(error);
        },
    });
});

$("#buttonSubmit").click(function () {
    alert("发送成功！")
    $.ajax({
        url: "http://127.0.0.1:9000/transactions",
        type: "GET",
        data: {
            "account1Address": $("#inputSenderAddress").val(),
            "account2Address": $("#inputReceiveAddress").val(),
            "money": $("#inputAmount").val()
        },
        success: function () {
            alert("发送成功！")
        }
    })
})



$("#nav-contact-tab").click(function () {
    $.ajax({
        url: "http://127.0.0.1:8080/history",
        type: "GET",
        dataType: 'json',
        success: function (response) {
            console.log(response);
            $('#List').empty();
            // 创建表格元素
            var table = $('<table>');

            // 创建表头行
            var headerRow = $('<tr>');
            headerRow.append($('<th>').text('Sender Address'));
            headerRow.append($('<th>').text('Receive Address'));
            headerRow.append($('<th>').text('Value'));
            table.append(headerRow);

            // 遍历响应数据并创建表格行
            $.each(response, function (index, item) {
                var row = $('<tr>');

                // 创建单元格并设置值
                var senderAddressCell = $('<td>').text(item.senderAddress);
                var receiveAddressCell = $('<td>').text(item.receiveAddress);
                var valueCell = $('<td>').text(item.value);

                // 将单元格添加到表格行
                row.append(senderAddressCell);
                row.append(receiveAddressCell);
                row.append(valueCell);

                // 将表格行添加到表格
                table.append(row);
            });
            // 将表格添加到页面中
            $('#List').append(table);
        }
    })
})

$("#trade").click(function () {
    $.ajax({
        url: "http://127.0.0.1:8080/history",
        type: "GET",
        dataType: 'json',
        success: function (response) {
            console.log(response);
            $('#List').empty();
            // 创建表格元素
            var table = $('<table>');

            // 创建表头行
            var headerRow = $('<tr>');
            headerRow.append($('<th>').text('Sender Address'));
            headerRow.append($('<th>').text('Receive Address'));
            headerRow.append($('<th>').text('Value'));
            table.append(headerRow);

            // 遍历响应数据并创建表格行
            $.each(response, function (index, item) {
                var row = $('<tr>');

                // 创建单元格并设置值
                var senderAddressCell = $('<td>').text(item.senderAddress);
                var receiveAddressCell = $('<td>').text(item.receiveAddress);
                var valueCell = $('<td>').text(item.value);

                // 将单元格添加到表格行
                row.append(senderAddressCell);
                row.append(receiveAddressCell);
                row.append(valueCell);

                // 将表格行添加到表格
                table.append(row);
            });
            // 将表格添加到页面中
            $('#List').append(table);
        }
    })
})

$("#nav-home-tab1").click(function () {
    window.location = "index0_bootstrap.html";
})

$("#nav-profile-tab1").click(function () {
    window.location = "index1_bootstrap.html";
})
$("#nav-contact-tab1").click(function () {
    window.location = "index2_bootstrap.html";
})



