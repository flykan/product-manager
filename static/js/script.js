let dataTable;

$(document).ready(function() {
    initializeDataTable();

    // 表单提交事件
    $('#productForm').on('submit', function(e) {
        e.preventDefault();
        createProduct();
    });
});

function initializeDataTable() {
    dataTable = $('#productsTable').DataTable({
        processing: true,
        serverSide: true,
        ajax: {
            url: '/api/products',
            type: 'GET',
            data: function(d) {
                return {
                    draw: d.draw,
                    page: (d.start / d.length) + 1,
                    limit: d.length,
                    'search[value]': d.search.value
                };
            },
            dataFilter: function(data) {
                let json = JSON.parse(data);
                return JSON.stringify({
                    draw: json.draw,
                    recordsTotal: json.recordsTotal,
                    recordsFiltered: json.recordsFiltered,
                    data: json.data
                });
            }
        },
        columns: [
            { data: 'ID' },
            { data: 'name' },
            { data: 'description' },
            {
                data: 'price',
                render: function(data) {
                    return '¥' + parseFloat(data).toFixed(2);
                }
            },
            { data: 'stock' },
            { data: 'category' },
            {
                data: 'CreatedAt',
                render: function(data) {
                    return new Date(data).toLocaleString();
                }
            },
            {
                data: null,
                render: function(data) {
                    return `
                        <div class="action-buttons">
                            <button class="btn btn-sm btn-warning" onclick="openEditModal(${data.ID})">编辑</button>
                            <button class="btn btn-sm btn-danger" onclick="deleteProduct(${data.ID})">删除</button>
                        </div>
                    `;
                },
                orderable: false
            }
        ]
    });
}

function createProduct() {
    const productData = {
        name: $('#name').val(),
        description: $('#description').val(),
        price: parseFloat($('#price').val()),
        stock: parseInt($('#stock').val()),
        category: $('#category').val()
    };

    $.ajax({
        url: '/api/products',
        type: 'POST',
        contentType: 'application/json',
        data: JSON.stringify(productData),
        success: function(response) {
            alert('商品添加成功！');
            clearForm();
            dataTable.ajax.reload();
        },
        error: function(xhr) {
            alert('添加商品失败: ' + (xhr.responseJSON?.error || '未知错误'));
        }
    });
}

function openEditModal(id) {
    $.ajax({
        url: `/api/products/${id}`,
        type: 'GET',
        success: function(product) {
            $('#editId').val(product.ID);
            $('#editName').val(product.name);
            $('#editDescription').val(product.description);
            $('#editPrice').val(product.price);
            $('#editStock').val(product.stock);
            $('#editCategory').val(product.category);

            $('#editModal').modal('show');
        },
        error: function() {
            alert('获取商品信息失败');
        }
    });
}

function updateProduct() {
    const productData = {
        name: $('#editName').val(),
        description: $('#editDescription').val(),
        price: parseFloat($('#editPrice').val()),
        stock: parseInt($('#editStock').val()),
        category: $('#editCategory').val()
    };

    const id = $('#editId').val();

    $.ajax({
        url: `/api/products/${id}`,
        type: 'PUT',
        contentType: 'application/json',
        data: JSON.stringify(productData),
        success: function(response) {
            alert('商品更新成功！');
            $('#editModal').modal('hide');
            dataTable.ajax.reload();
        },
        error: function(xhr) {
            alert('更新商品失败: ' + (xhr.responseJSON?.error || '未知错误'));
        }
    });
}

function deleteProduct(id) {
    if (!confirm('确定要删除这个商品吗？')) {
        return;
    }

    $.ajax({
        url: `/api/products/${id}`,
        type: 'DELETE',
        success: function(response) {
            alert('商品删除成功！');
            dataTable.ajax.reload();
        },
        error: function(xhr) {
            alert('删除商品失败: ' + (xhr.responseJSON?.error || '未知错误'));
        }
    });
}

function clearForm() {
    $('#productForm')[0].reset();
}