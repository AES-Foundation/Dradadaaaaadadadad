let settings = () => {
    $("#hideInDonaters").click((e) => {
        $("#hideInDonatersCheckbox #loading").removeClass('d-none');
        $.ajax({
            url: `/a/settings?hideInTopDonaters=${e.target.checked}`,
            success: (r) => {
                $("#hideInDonatersCheckbox #loading").addClass('d-none');
            },
        })
    })
    
    let t = -1;
    
    $("#displayName").on('keyup', (e) => {
        if (t > 0) {
            clearTimeout(t);
        }
        t = setTimeout(() => {
            let val = $("#displayName").val();
            if (val.length == 0) {
                return;
            }
    
            $("#displayNameField #loading").removeClass('d-none');
            $.ajax({
                url: `/a/settings?displayName=${val}`,
                success: (r) => {
                    $("#displayNameField #loading").addClass('d-none');
                    $("#profileDisplayName").text(val)
                },
            })
        }, 500);
    })
    
    $("#clearHistory").click(() => {
        $("#clearHistoryField #loading").removeClass('d-none');
        $.ajax({
            url: `/a/settings/clearHistory`,
            success: (r) => {
                $("#clearHistoryField #loading").addClass('d-none');
            },
        });
    });
}
let avatar = () => {
    $("#generate").on("submit", (e) => {
        e.preventDefault();
        $('#error').text('');
        
        if ($('#nickname').val() == "") {
            alert("Поле не может быть пустоё!")
            return;
        }

        $('#head').addClass('loading');

        $.ajax({
            url: `/a/generate?data=${btoa($('#nickname').val())}&params=${ $("#voxel").get()[0].checked ? "voxel" : "default" }`,
            success: (r) => {
                $('#head').removeClass('loading');
                if (r.error != null) {
                    switch (r.error) {
                        case 'invalid source':
                            $('#error').text("Неверный запрос");
                            break;
                        default:
                            console.log(r.error);
                            break;
                    }
                    return
                }
                $('#head img').attr("src", r.url)
                let headForList = $(`<div class="m-3 ms-0">
                        <img src="${r.url}" width="64">
                    </div>`);
                $('#headsList').prepend(headForList);
            },
        })

        let phrases = [
            'Жить на что-то надо 😖',
            'Мы для вас старались 🥺',
            'У меня лапки 🐾',
            'Хочу немного деняг 💵',
            'Сколько это дошираков? 🤔',
            'На новую коробку 📦',
        ]
        let phrase = phrases[Math.round(Math.random() * (phrases.length - 1))]
        let modal = $(`<div class="picraft-modal">
            <div class="card picraft-modal-card">
                <div class="card-header d-flex align-items-center">
                    Донатик
                </div>
                <div class="m-3">
                    <div class="mb-3">${phrase}</div>
                    <div class="d-flex">
                        <button id="donate" class="picraft-btn picraft-btn-white flex-grow-1 me-3"><img src="/static/qiwi_logo_rgb_small.png" height="40"></button>
                        <button id="cancel" class="picraft-btn picraft-btn flex-shrink-0">Нет</button>
                    </div>
                </div>
            </div>
        </div>`);
        modal.find('#donate').click(() => {
            window.open('/donate', 'donate', 'width=870,height=660');
            modal.remove();
        });
        modal.find('#cancel').click(() => {
            modal.remove();
        });
        $(document.body).append(modal)
    })

    $("#download").click(() => {
        $(`<a href="${$('#head img').attr("src")}" download></a>`).each((_, e) => {
            e.click();
        })
    });

    $("#voxelCheckbox").click(() => {
        if ($("#voxel").get()[0].disabled) {
            alert("Данная функция доступно от суммы 50 рублей")
        }
    });
}

let join = (array, callback) => {
    let a = ''
    array.forEach(element => {
        a += callback(element)
    });
    return a;
}

let Threehead = {
    'move': {
        'settings': () => {
            $('#content').addClass('picraft-blur');
            $.ajax({
                url: "/settings?json=true",
                success: (r) => {
                    let parent = $('#content').parent()
                    $('#content').remove()
                    let newContent = $(`<div class="card w-100 mt-3" id="content">
                                        <div class="card-header d-flex align-items-center">
                                            Настройки
                                        </div>
                                        <form class="p-3">
                                            <div class="d-flex align-items-center mb-3" id="displayNameField">
                                                <label for="displayName" class="me-3">Имя</label>
                                                <input type="text" class="form-control w-75" value="${r.user.displayName}" id="displayName">
                                                <div class="spinner-border text-primary ms-3 d-none" role="status" id="loading" style="width: 1.2rem; height: 1.2rem;"></div>
                                            </div>
                                            <div class="d-flex align-items-center mb-3" id="hideInDonatersCheckbox">
                                                <input type="checkbox" id="hideInDonaters" class="picraft-switch me-1" ${r.user.isHiddenInTopDonaters ? 'checked' : ''}>
                                                <label for="voxel">Скрыть меня из список топ донатеров</label>
                                                <div class="spinner-border text-primary ms-3 d-none" role="status" id="loading" style="width: 1.2rem; height: 1.2rem;"></div>
                                            </div>
                                            <div class="d-flex align-items-center" id="clearHistoryField">
                                                <button type="button" class="picraft-btn picraft-btn-primary flex-shrink-0" id="clearHistory">Очистить историю</button>
                                                <div class="spinner-border text-primary ms-3 d-none" role="status" id="loading" style="width: 1.2rem; height: 1.2rem;"></div>
                                            </div>
                                        </form>
                                    </div>`
                    );
                    parent.append(newContent);
                    settings();
                    history.pushState('', '', '/settings');
                }
            })
        },
        'avatar': () => {
            $('#content').addClass('picraft-blur');
            $.ajax({
                url: "/avatar?json=true",
                success: (r) => {
                    let parent = $('#content').parent()
                    $('#content').remove()
                    let newContent = $(`
                           <div id="content">
                                <div class="d-flex mt-3">
                                    <div class="flex-grow-1 me-3 card w-100">
                                        <div class="card-header d-flex align-items-center">
                                            <img src="/static/icon.png" class="me-3" width="64">Генератор
                                        </div>
                                        <div class="m-3 d-flex flex-row">
                                            <div class="me-3 position-relative picraft-head" id="head">
                                                <img src="/static/icon.png">
                                                <div class="position-absolute top-50 start-50 spinner-border text-primary picraft-spinner" role="status"></div>
                                            </div>
                                            <form class="flex-grow-1 flex-shrink-1" id="generate">
                                                <input type="text" class="form-control mb-3" placeholder="Ник Майнкрафт" id="nickname">
                                                <div class="d-flex">
                                                    <button type="submit" class="picraft-btn picraft-btn-primary flex-grow-1 flex-shrink-1 me-3">Сгенерировать</button>
                                                    <button type="button" class="picraft-btn flex-shrink-0" id="download" style="flex-grow: 0.5;">Скачать</button>
                                                </div>
                                                <div class="mt-3 d-flex align-items-center" id="voxelCheckbox">
                                                    <input type="checkbox" id="voxel" class="picraft-switch me-1">
                                                    <label for="voxel">Режим Вокселя</label>
                                                </div>
                                                <div class="mt-3 picraft-error" id="error">
                                                </div>
                                            </form>
                                        </div>
                                    </div>
                                    <div style="width: 22rem;">
                                        <div class="card w-100 h-100">
                                            <div class="card-header d-flex align-items-center">
                                                <img src="/static/donate.png" class="me-3" width="64">ТОП ДОНАТЕРЫ
                                            </div>
                                            <div class="pt-3" style="overflow-y: auto; max-height: 14rem;">
                                                ${join(r.topDonaters, (e) => `
                                                <div class="m-3 mt-0">
                                                    <b>${e.amount} рублей</b> <span>${e.displayName}</span>
                                                </div>
                                                `)}
                                            </div>
                                        </div>
                                    </div>
                                </div>
                                <div class="card w-100 mt-3">
                                    <div class="card-header d-flex align-items-center">
                                        Последние твои запросы
                                    </div>
                                    <div class="ps-3 d-flex" style="overflow-x: auto;" id="headsList">
                                        ${join(r.latestHeads, (e) => `
                                        <div class="m-3 ms-0">
                                            <img src="${e}" width="64">
                                        </div>
                                        `)}
                                    </div>
                                </div>
                            </div>
                        </div>`
                    );
                    parent.append(newContent);
                    avatar();
                    history.pushState('', '', '/avatar');
                }
            })
        }
    }
}