create table "user"
(
    id                varchar(36)                                                not null
        primary key,
    email             varchar(100)                                               not null
        unique,
    password          varchar(255)                                               not null,
    nickname          varchar(50)              default ''::character varying,
    avatar            varchar(255)             default ''::character varying,
    status            integer                  default 1,
    last_login_at     timestamp,
    created_at        timestamp with time zone default CURRENT_TIMESTAMP,
    updated_at        timestamp with time zone default CURRENT_TIMESTAMP,
    deleted_at        timestamp with time zone,
    tags              varchar(255)             default ''::character varying,
    identity          smallint                 default 1                         not null,
    max_project_count integer                  default 0                         not null,
    max_team_count    integer                  default 0                         not null,
    role              varchar(10)              default 'user'::character varying not null
);

comment on table "user" is '用户表';

comment on column "user".id is '用户唯一标识符(UUID)';

comment on column "user".email is '邮箱地址，用于登录和通知，全局唯一';

comment on column "user".password is '加密后的密码';

comment on column "user".nickname is '用户昵称，用于显示';

comment on column "user".avatar is '头像URL地址';

comment on column "user".status is '用户状态：1=正常 2=禁用';

comment on column "user".last_login_at is '最后登录时间';

comment on column "user".created_at is '账户创建时间';

comment on column "user".updated_at is '账户信息最后更新时间';

comment on column "user".deleted_at is '软删除时间';

comment on column "user".identity is '权益类型: 1=基础版 2=专业版 3=旗舰版';

comment on column "user".role is '角色: admin=管理员 user=用户';

alter table "user"
    owner to postgres;

create index idx_user_email
    on "user" (email);

comment on index idx_user_email is '邮箱索引，用于快速登录验证';

create index idx_user_status
    on "user" (status);

comment on index idx_user_status is '用户状态索引，用于筛选活跃用户';

create index idx_user_created_at
    on "user" (created_at);

comment on index idx_user_created_at is '创建时间索引，用于时间范围查询';

create index idx_user_deleted_at
    on "user" (deleted_at);

comment on index idx_user_deleted_at is '删除时间索引，用于软删除查询';

create index idx_user_identity
    on "user" (identity);

create table invite_code
(
    id         varchar(36) not null
        primary key,
    user_id    varchar(36) not null,
    value      varchar(20) not null
        unique,
    used_by    varchar(256),
    is_used    boolean                  default false,
    created_at timestamp with time zone default CURRENT_TIMESTAMP,
    updated_at timestamp with time zone default CURRENT_TIMESTAMP,
    deleted_at timestamp with time zone
);

comment on table invite_code is '邀请码表';

comment on column invite_code.id is '主键';

comment on column invite_code.user_id is '创建邀请码的用户ID';

comment on column invite_code.value is '邀请码';

comment on column invite_code.used_by is '使用邀请码的用户邮箱';

alter table invite_code
    owner to postgres;

create index idx_invite_code_value
    on invite_code (value);

create index idx_invite_code_deleted_at
    on invite_code (deleted_at);

create table file
(
    id         varchar(36)                                        not null
        primary key,
    user_id    varchar(36)                                        not null,
    filename   varchar(255)             default ''::character varying,
    filesize   bigint                   default 0,
    filetype   varchar(100)             default ''::character varying,
    oss        varchar(20)              default ''::character varying,
    path       varchar(255)             default ''::character varying,
    remark     varchar(255)             default ''::character varying,
    created_at timestamp with time zone default CURRENT_TIMESTAMP not null,
    updated_at timestamp with time zone default CURRENT_TIMESTAMP not null,
    deleted_at timestamp with time zone
);

comment on table file is '文件表';

comment on column file.id is '主键id';

comment on column file.user_id is '用户ID';

comment on column file.filename is '文件名';

comment on column file.filesize is '文件大小，单位字节';

comment on column file.filetype is '文件类型';

comment on column file.oss is '对象存储: aliyun/minio';

comment on column file.path is '文件存储路径';

comment on column file.remark is '备注';

comment on column file.created_at is '创建时间';

comment on column file.updated_at is '更新时间';

comment on column file.deleted_at is '删除时间';

alter table file
    owner to postgres;

create index idx_file_info_deleted_at
    on file (deleted_at);

create index idx_file_oss
    on file (oss);

create table website_config
(
    id         varchar(36) not null
        primary key,
    schema     jsonb       not null,
    created_at timestamp with time zone default CURRENT_TIMESTAMP,
    updated_at timestamp with time zone default CURRENT_TIMESTAMP,
    deleted_at timestamp with time zone
);

alter table website_config
    owner to postgres;

create table team
(
    id          varchar(36)                                            not null
        primary key,
    user_id     varchar(36)                                            not null,
    name        varchar(255)             default ''::character varying not null,
    description varchar(1024)            default ''::character varying not null,
    logo        varchar(255)             default ''::character varying not null,
    created_at  timestamp with time zone default CURRENT_TIMESTAMP     not null,
    updated_at  timestamp with time zone default CURRENT_TIMESTAMP     not null,
    deleted_at  timestamp with time zone
);

comment on column team.user_id is '创建人id';

comment on column team.name is '团队名称';

comment on column team.description is '团队描述';

comment on column team.logo is '团队logo';

alter table team
    owner to postgres;

create index idx_team_deleted_at
    on team (deleted_at);

create index idx_team_user_id
    on team (user_id);

create table user_team
(
    id         varchar(36)                                                  not null
        primary key,
    user_id    varchar(36)                                                  not null,
    team_id    varchar(36)                                                  not null,
    identity   varchar(10)              default 'member'::character varying not null,
    created_at timestamp with time zone default CURRENT_TIMESTAMP           not null,
    updated_at timestamp with time zone default CURRENT_TIMESTAMP           not null,
    deleted_at timestamp with time zone
);

comment on column user_team.user_id is '用户id';

comment on column user_team.team_id is '团队id';

comment on column user_team.identity is '身份: creator=创建人 manager=创建人 member=成员';

alter table user_team
    owner to postgres;

create index idx_user_team_user_id
    on user_team (user_id);

create index idx_user_team_team_id
    on user_team (team_id);

create index idx_user_team_deleted_at
    on user_team (deleted_at);

create table mv_project
(
    id                   varchar(36)                            not null
        primary key,
    name                 varchar(255)                           not null,
    description          text                     default ''::text,
    user_id              varchar(36)                            not null,
    created_at           timestamp with time zone default CURRENT_TIMESTAMP,
    updated_at           timestamp with time zone default CURRENT_TIMESTAMP,
    deleted_at           timestamp with time zone,
    enable_advanced_perm boolean                  default false not null
);

comment on column mv_project.enable_advanced_perm is '是否开启高级权限';

alter table mv_project
    owner to postgres;

create index idx_mv_project_user_id
    on mv_project (user_id);

create table mv_folder
(
    id         varchar(36)                                            not null
        primary key,
    project_id varchar(36)                                            not null,
    name       varchar(255)                                           not null,
    parent_id  varchar(36)              default ''::character varying not null,
    created_at timestamp with time zone default CURRENT_TIMESTAMP,
    updated_at timestamp with time zone default CURRENT_TIMESTAMP,
    deleted_at timestamp with time zone
);

alter table mv_folder
    owner to postgres;

create index idx_mv_folder_project_id
    on mv_folder (project_id);

create index idx_mv_folder_parent_id
    on mv_folder (parent_id);

create table mv_table_schema
(
    id          varchar(36)                                            not null
        primary key,
    folder_id   varchar(36)              default ''::character varying not null,
    name        varchar(255)                                           not null,
    version     integer                  default 1                     not null,
    created_by  varchar(36)                                            not null,
    updated_by  varchar(36),
    description text                     default ''::text,
    config      jsonb                    default '{}'::jsonb,
    stats       jsonb                    default '{}'::jsonb,
    created_at  timestamp with time zone default CURRENT_TIMESTAMP,
    updated_at  timestamp with time zone default CURRENT_TIMESTAMP,
    deleted_at  timestamp with time zone,
    row_name    varchar(50)              default ''::character varying not null,
    project_id  text                     default ''::text              not null
);

comment on column mv_table_schema.row_name is '每一行的名称';

alter table mv_table_schema
    owner to postgres;

create index idx_mv_table_schema_folder_id
    on mv_table_schema (folder_id);

create index idx_mv_table_schema_name
    on mv_table_schema (name);

create index idx_mv_table_schema_created_by
    on mv_table_schema (created_by);

create table mv_field
(
    id              varchar(36)                                                      not null
        primary key,
    table_schema_id text                                                             not null,
    title           varchar(255)                                                     not null,
    type            varchar(50)                                                      not null,
    config          jsonb                    default '{}'::jsonb,
    order_index     varchar(36)              default 'mmmmmmmmmm'::character varying not null,
    created_at      timestamp with time zone default CURRENT_TIMESTAMP,
    updated_at      timestamp with time zone default CURRENT_TIMESTAMP,
    deleted_at      timestamp with time zone,
    is_primary_key  boolean                  default false                           not null
);

comment on column mv_field.title is '字段标题';

comment on column mv_field.type is '字段类型';

comment on column mv_field.config is '字段配置';

comment on column mv_field.is_primary_key is '是否为主键';

alter table mv_field
    owner to postgres;

create index idx_mv_field_table_id
    on mv_field (table_schema_id);

create table mv_record
(
    id              varchar(36)                                                      not null
        primary key,
    table_schema_id text                                                             not null,
    created_by      varchar(36)                                                      not null,
    row_data        jsonb                    default '{}'::jsonb,
    created_at      timestamp with time zone default CURRENT_TIMESTAMP,
    updated_at      timestamp with time zone default CURRENT_TIMESTAMP,
    deleted_at      timestamp with time zone,
    order_index     varchar(36)              default 'mmmmmmmmmm'::character varying not null,
    updated_by      text
);

comment on column mv_record.updated_by is '更新人id';

alter table mv_record
    owner to postgres;

create index idx_mv_record_created_by
    on mv_record (created_by);

create index idx_mv_record_values_gin
    on mv_record using gin (row_data jsonb_path_ops);

create index idx_mv_record_row_data_gin
    on mv_record using gin (row_data);

create index idx_mv_record_table_id
    on mv_record (table_schema_id);

create table mv_project_perm
(
    id         varchar(36)                                                  not null
        constraint mv_project_permission_pkey
            primary key,
    target_id  varchar(36)                                                  not null,
    project_id varchar(36)                                                  not null,
    role       varchar(20)              default 'reader'::character varying not null,
    created_at timestamp with time zone default CURRENT_TIMESTAMP,
    updated_at timestamp with time zone default CURRENT_TIMESTAMP,
    deleted_at timestamp with time zone,
    target     varchar(20)              default 'user'::character varying   not null
);

comment on table mv_project_perm is '项目用户角色表';

comment on column mv_project_perm.target_id is '用户id/团队id';

comment on column mv_project_perm.role is '项目角色: 所有者=owner 管理员=admin 编辑者=editor 阅读者=reader';

comment on column mv_project_perm.target is 'user/team';

alter table mv_project_perm
    owner to postgres;

create index mv_project_perm_project_id_target_target_id_index
    on mv_project_perm (project_id, target, target_id)
    where (deleted_at IS NULL);

create table mv_project_state
(
    id              varchar(36)                                                not null
        primary key,
    project_id      varchar(36)                                                not null,
    share_range     integer                  default 1                         not null,
    team_action     varchar(20)              default 'read'::character varying not null,
    access_password varchar(255)             default ''::character varying     not null,
    created_at      timestamp with time zone default CURRENT_TIMESTAMP,
    updated_at      timestamp with time zone default CURRENT_TIMESTAMP,
    deleted_at      timestamp with time zone
);

comment on column mv_project_state.share_range is '分享范围: 1=仅协作者可见, 2=互联网公开, 3=团队内公开';

comment on column mv_project_state.team_action is '团队内权限(分享范围为团队内公开时该字段有意义): read=团队内可查看, edit=团队内可编辑';

comment on column mv_project_state.access_password is '访问密码(分享范围为互联网公开时该字段有意义)';

alter table mv_project_state
    owner to postgres;

create index idx_resource_state_resource_id
    on mv_project_state (project_id);

create index idx_resource_state_deleted_at
    on mv_project_state (deleted_at);

create table mv_project_advanced_perm
(
    id                varchar(36)                                                not null
        constraint mv_advanced_permission_pkey
            primary key,
    table_schema_id   varchar(36)                                                not null,
    role              varchar(36)                                                not null,
    data_action       varchar(20)              default 'none'::character varying not null,
    can_add           boolean                  default false                     not null,
    can_delete        boolean                  default false                     not null,
    operate_range     varchar(20)              default 'all'::character varying  not null,
    field_access      varchar(20)              default 'all'::character varying  not null,
    custom_field_perm jsonb                    default '[]'::jsonb               not null,
    view_access       varchar(20)              default 'all'::character varying  not null,
    can_check_views   jsonb                    default '[]'::jsonb               not null,
    created_at        timestamp with time zone default CURRENT_TIMESTAMP,
    updated_at        timestamp with time zone default CURRENT_TIMESTAMP,
    deleted_at        timestamp with time zone,
    can_operate_view  boolean                  default false                     not null
);

comment on table mv_project_advanced_perm is '项目高级权限表';

comment on column mv_project_advanced_perm.role is '角色: 所有者=owner 管理员=admin 编辑者=editor 阅读者=reader';

comment on column mv_project_advanced_perm.data_action is '数据表权限: 可管理=manage 可编辑=edit 可查看=read 无权限=none';

comment on column mv_project_advanced_perm.operate_range is '操作范围: 所有记录=all 仅自己创建的记录=own';

comment on column mv_project_advanced_perm.field_access is '字段访问权限: 所有字段=all 自定义权限=custom';

comment on column mv_project_advanced_perm.custom_field_perm is '自定义字段权限: [{ "fieldId": "字段ID", "canRead": true, "canAdd": true, "canEdit": true }]';

comment on column mv_project_advanced_perm.view_access is '视图访问权限: 所有视图=all 自定义视图=custom';

comment on column mv_project_advanced_perm.can_check_views is '可查看的视图: ["视图ID1", "视图ID2"]';

comment on column mv_project_advanced_perm.can_operate_view is '是否可新增、修改、删除视图';

alter table mv_project_advanced_perm
    owner to postgres;

create index idx_mv_advanced_permission_table_schema_id
    on mv_project_advanced_perm (table_schema_id);

create index idx_mv_advanced_permission_role
    on mv_project_advanced_perm (role);

create index idx_mv_advanced_permission_table_schema_id_role
    on mv_project_advanced_perm (table_schema_id, role);

create index idx_mv_advanced_permission_deleted_at
    on mv_project_advanced_perm (deleted_at);

create table mv_project_favorite
(
    id         varchar(36) not null
        primary key,
    project_id varchar(36) not null,
    user_id    varchar(36) not null,
    created_at timestamp with time zone default CURRENT_TIMESTAMP,
    updated_at timestamp with time zone default CURRENT_TIMESTAMP,
    deleted_at timestamp with time zone
);

alter table mv_project_favorite
    owner to postgres;

create index idx_mv_project_favorite_user_id
    on mv_project_favorite (user_id);

create index idx_mv_project_favorite_deleted_at
    on mv_project_favorite (deleted_at);

create index mv_project_favorite_user_id_project_id_index
    on mv_project_favorite (user_id, project_id);

create table mv_project_recent
(
    id               varchar(36) not null
        primary key,
    project_id       varchar(36) not null,
    user_id          varchar(36) not null,
    last_accessed_at timestamp with time zone default CURRENT_TIMESTAMP,
    created_at       timestamp with time zone default CURRENT_TIMESTAMP,
    updated_at       timestamp with time zone default CURRENT_TIMESTAMP,
    deleted_at       timestamp with time zone
);

alter table mv_project_recent
    owner to postgres;

create index idx_mv_project_recent_user_id
    on mv_project_recent (user_id);

create index idx_mv_project_recent_deleted_at
    on mv_project_recent (deleted_at);

create unique index idx_mv_project_recent_user_id_project_id
    on mv_project_recent (user_id, project_id);

create table mv_view
(
    id              varchar(36)                                         not null
        primary key,
    table_schema_id varchar(36)                                         not null,
    type            text                     default 'table'::text      not null,
    created_at      timestamp with time zone default CURRENT_TIMESTAMP,
    updated_at      timestamp with time zone default CURRENT_TIMESTAMP,
    deleted_at      timestamp with time zone,
    name            text                     default ''::text           not null,
    description     text                     default ''::text           not null,
    order_index     text                     default 'mmmmmmmmmm'::text not null
);

comment on table mv_view is '数据视图';

comment on column mv_view.table_schema_id is '数据表格id';

comment on column mv_view.type is '视图类型: table=表格视图, form=表单视图';

comment on column mv_view.name is '视图名称';

alter table mv_view
    owner to postgres;

create index idx_mv_view_deleted_at
    on mv_view (deleted_at);

create index idx_mv_view_order_index
    on mv_view (order_index);

create index idx_mv_view_active_deleted_at
    on mv_view (deleted_at);

create table mv_view_table
(
    id            varchar(36) not null
        primary key,
    filter_config jsonb                    default '[]'::jsonb,
    group_config  jsonb                    default '{}'::jsonb,
    sort_config   jsonb                    default '[]'::jsonb,
    row_height    integer                  default 1,
    color_config  jsonb                    default '[]'::jsonb,
    created_at    timestamp with time zone default CURRENT_TIMESTAMP,
    updated_at    timestamp with time zone default CURRENT_TIMESTAMP,
    deleted_at    timestamp with time zone,
    view_id       varchar(36) not null
);

comment on table mv_view_table is '表格视图';

comment on column mv_view_table.filter_config is '筛选条件: [{fieldId: "字段id", operator: "equal=等于/notEqual=不等于/contains=包含/notContains=不包含/null=为空/notNull=不为空/before=早于/after=晚于/beforeOrEqual=早于或等于/afterOrEqual=晚于或等于", value: "值"}]';

comment on column mv_view_table.group_config is '分组配置: {fieldId: "", order: "asc/desc"}';

comment on column mv_view_table.sort_config is '排序配置: [{fieldId: "", order: "asc/desc"}]';

comment on column mv_view_table.row_height is '行高: 1=常规, 2=中等, 3=高, 4=超高';

comment on column mv_view_table.color_config is '颜色配置: [{color: "#cccccc", scope: "cell/row/rowHeader/column", fieldId: "字段id", operator: "equal=等于/notEqual=不等于/greaterThan=大于/lessThan=小于/greaterThanOrEqual=大于等于/lessThanOrEqual=小于等于/contains=包含/notContains=不包含/null=为空/notNull=不为空/before=早于/after=晚于/beforeOrEqual=早于或等于/afterOrEqual=晚于或等于", value: "值"}]';

alter table mv_view_table
    owner to postgres;

create index idx_mv_view_table_deleted_at
    on mv_view_table (deleted_at);

create table mv_view_form
(
    id                       varchar(36)                                     not null
        primary key,
    name                     text                     default ''::text       not null,
    description              text                     default ''::text       not null,
    cover                    text                     default ''::text       not null,
    layout                   text                     default 'center'::text not null,
    stats                    jsonb                    default '{}'::jsonb    not null,
    enable_sharing           boolean                  default false          not null,
    enable_anonymous         boolean                  default false          not null,
    filter                   text                     default 'all'::text    not null,
    filter_config            jsonb                    default '[]'::jsonb    not null,
    enable_no_login          boolean                  default false          not null,
    enable_limit_submit      boolean                  default false          not null,
    limit_submit_type        text                     default 'once'::text   not null,
    enable_limit_collect     boolean                  default false          not null,
    limit_collect_count      integer                  default 0              not null,
    enable_cycle_remind      boolean                  default false          not null,
    cycle_remind_config      jsonb                    default '{}'::jsonb    not null,
    enable_edit_after_submit boolean                  default false          not null,
    created_at               timestamp with time zone default CURRENT_TIMESTAMP,
    updated_at               timestamp with time zone default CURRENT_TIMESTAMP,
    deleted_at               timestamp with time zone,
    view_id                  varchar(36)                                     not null,
    config                   jsonb                    default '[]'::jsonb    not null
);

comment on table mv_view_form is '表单视图';

comment on column mv_view_form.name is '表单名称';

comment on column mv_view_form.description is '表单描述';

comment on column mv_view_form.cover is '表单封面图片';

comment on column mv_view_form.layout is '表单布局: center/left';

comment on column mv_view_form.stats is '表单统计信息';

comment on column mv_view_form.enable_sharing is '是否开启分享';

comment on column mv_view_form.enable_anonymous is '是否允许匿名填写';

comment on column mv_view_form.filter is '表单填写权限: all=所有人, specific=指定用户/团队';

comment on column mv_view_form.filter_config is '表单填写权限配置: [{targetId: "用户id/团队id", target: "user/team"}]';

comment on column mv_view_form.enable_no_login is '是否允许未登录用户填写';

comment on column mv_view_form.enable_limit_submit is '是否开启填写次数限制';

comment on column mv_view_form.limit_submit_type is '填写次数限制类型: once=只填写一次, perDay=每天填写一次, perCycle=每个周期填写一次(开启按周期定时提醒时生效)';

comment on column mv_view_form.enable_limit_collect is '是否开启收集上限限制';

comment on column mv_view_form.limit_collect_count is '收集上限数量';

comment on column mv_view_form.enable_cycle_remind is '(暂不使用)是否开启按周期定时提醒';

comment on column mv_view_form.cycle_remind_config is '(暂不使用)周期定时提醒配置: {type: "daily/weekly/monthly", submitDate: ["monday", "tuesday", "wednesday", "thursday", "friday", "saturday", "sunday"], startTime: "", endTime: ""}';

comment on column mv_view_form.enable_edit_after_submit is '是否允许提交后编辑';

comment on column mv_view_form.view_id is '视图id';

comment on column mv_view_form.config is '表单项配置: [{fieldId: "", description: "", defaultValue: "", required: false, isHide: false}]';

alter table mv_view_form
    owner to postgres;

create index idx_mv_view_form_deleted_at
    on mv_view_form (deleted_at);

create table mv_form_submit
(
    id           varchar(36) not null
        primary key,
    user_id      varchar(36) not null,
    view_form_id varchar(36) not null,
    record_id    varchar(36) not null,
    created_at   timestamp with time zone default CURRENT_TIMESTAMP,
    updated_at   timestamp with time zone default CURRENT_TIMESTAMP,
    deleted_at   timestamp with time zone
);

comment on table mv_form_submit is '表单提交记录表';

comment on column mv_form_submit.user_id is '表单提交人id';

comment on column mv_form_submit.view_form_id is '表单视图id';

comment on column mv_form_submit.record_id is '表单提交后对应的表格数据记录id';

alter table mv_form_submit
    owner to postgres;

create index idx_mv_form_submit_view_form_id
    on mv_form_submit (view_form_id);

create index idx_mv_form_submit_user_id
    on mv_form_submit (user_id);

create index idx_mv_form_submit_deleted_at
    on mv_form_submit (deleted_at);

create table mv_template_project
(
    id          text                                               not null
        primary key,
    user_id     text                                               not null,
    name        text                     default ''::text          not null,
    description text                     default ''::text,
    cover       text                     default ''::text,
    tags        text[]                   default '{}'::text[],
    use_count   integer                  default 0                 not null,
    created_at  timestamp with time zone default CURRENT_TIMESTAMP not null,
    updated_at  timestamp with time zone default CURRENT_TIMESTAMP not null,
    deleted_at  timestamp with time zone
);

alter table mv_template_project
    owner to postgres;

create index idx_mv_template_project_user_id
    on mv_template_project (user_id);

create index idx_mv_template_project_created_at
    on mv_template_project (created_at);

create index idx_mv_template_project_deleted_at
    on mv_template_project (deleted_at);

create index idx_mv_template_project_tags_gin
    on mv_template_project using gin (tags);

create table mv_template_tag
(
    id          text                                                not null
        primary key,
    name        text                     default ''::text           not null,
    description text                     default ''::text,
    order_index text                     default 'mmmmmmmmmm'::text not null,
    created_at  timestamp with time zone default CURRENT_TIMESTAMP  not null,
    updated_at  timestamp with time zone default CURRENT_TIMESTAMP  not null,
    deleted_at  timestamp with time zone
);

alter table mv_template_tag
    owner to postgres;

create index idx_mv_template_tag_order_index
    on mv_template_tag (order_index);

create index idx_mv_template_tag_created_at
    on mv_template_tag (created_at);

create index idx_mv_template_tag_deleted_at
    on mv_template_tag (deleted_at);

create table mv_dashboard_chart
(
    id              text                                      not null
        primary key,
    dashboard_id    text                     default ''::text not null,
    table_schema_id text                     default ''::text not null,
    field1_id       text                     default ''::text not null,
    field2_id       text                     default ''::text not null,
    title           text                     default ''::text not null,
    type            text                     default ''::text not null,
    config          jsonb                    default '{}'::jsonb,
    created_at      timestamp with time zone default CURRENT_TIMESTAMP,
    updated_at      timestamp with time zone default CURRENT_TIMESTAMP,
    deleted_at      timestamp with time zone
);

alter table mv_dashboard_chart
    owner to postgres;

create index idx_mv_dashboard_chart_dashboard
    on mv_dashboard_chart (dashboard_id);

create table mv_dashboard
(
    id         text                                      not null
        primary key,
    project_id text                     default ''::text not null,
    folder_id  text                     default ''::text not null,
    name       text                     default ''::text not null,
    theme      text                     default ''::text not null,
    bg_type    text                     default ''::text not null,
    bg_color   text                     default ''::text not null,
    bg_image   text                     default ''::text not null,
    created_at timestamp with time zone default CURRENT_TIMESTAMP,
    updated_at timestamp with time zone default CURRENT_TIMESTAMP,
    deleted_at timestamp with time zone
);

comment on column mv_dashboard.bg_type is '背景颜色类型: followTheme/color/image';

alter table mv_dashboard
    owner to postgres;

create index idx_mv_dashboard_project
    on mv_dashboard (project_id);

create index idx_mv_dashboard_project_folder
    on mv_dashboard (project_id, folder_id);

create index idx_mv_dashboard_created_at
    on mv_dashboard (created_at);

create index idx_mv_dashboard_deleted_at
    on mv_dashboard (deleted_at);

create index idx_mv_dashboard_chart_created_at
    on mv_dashboard (created_at);

create index idx_mv_dashboard_chart_deleted_at
    on mv_dashboard (deleted_at);

create table mv_rich_text_content
(
    id         varchar(36)                               not null
        primary key,
    record_id  varchar(36)                               not null,
    field_id   varchar(36)                               not null,
    content    text                     default ''::text not null,
    created_at timestamp with time zone default CURRENT_TIMESTAMP,
    updated_at timestamp with time zone default CURRENT_TIMESTAMP,
    deleted_at timestamp with time zone
);

comment on table mv_rich_text_content is '富文本内容表';

comment on column mv_rich_text_content.record_id is '记录ID，关联 mv_record.id';

comment on column mv_rich_text_content.field_id is '字段ID，关联 mv_field.id';

comment on column mv_rich_text_content.content is '富文本内容（HTML格式）';

alter table mv_rich_text_content
    owner to postgres;

create unique index idx_mv_rich_text_content_record_field
    on mv_rich_text_content (record_id, field_id)
    where (deleted_at IS NULL);

create index idx_mv_rich_text_content_record_id
    on mv_rich_text_content (record_id);

create index idx_mv_rich_text_content_field_id
    on mv_rich_text_content (field_id);

create index idx_mv_rich_text_content_deleted_at
    on mv_rich_text_content (deleted_at);

create table mv_view_active
(
    id              text not null
        primary key,
    table_schema_id text not null,
    view_id         text not null,
    user_id         text not null,
    created_at      timestamp with time zone default CURRENT_TIMESTAMP,
    updated_at      timestamp with time zone default CURRENT_TIMESTAMP,
    deleted_at      timestamp with time zone
);

alter table mv_view_active
    owner to postgres;

create index idx_mv_table_schema_id_view_id_user_id
    on mv_view_active (table_schema_id, view_id, user_id);

create table mv_view_board
(
    id               text not null
        primary key,
    view_id          text not null,
    show_field_title boolean                  default false,
    filter_config    jsonb                    default '[]'::jsonb,
    group_config     jsonb                    default '{}'::jsonb,
    sort_config      jsonb                    default '[]'::jsonb,
    created_at       timestamp with time zone default CURRENT_TIMESTAMP,
    updated_at       timestamp with time zone default CURRENT_TIMESTAMP,
    deleted_at       timestamp with time zone
);

comment on table mv_view_board is '看板视图';

comment on column mv_view_board.show_field_title is '是否显示字段标题';

comment on column mv_view_board.filter_config is '筛选条件: [{fieldId: "字段id", operator: "equal=等于/notEqual=不等于/contains=包含/notContains=不包含/null=为空/notNull=不为空/before=早于/after=晚于/beforeOrEqual=早于或等于/afterOrEqual=晚于或等于", value: "值"}]';

comment on column mv_view_board.group_config is '分组配置: {fieldId: ""}';

comment on column mv_view_board.sort_config is '排序配置: [{fieldId: "", order: "asc/desc"}]';

alter table mv_view_board
    owner to postgres;

create index idx_mv_view_board_deleted_at
    on mv_view_board (deleted_at);

create table mv_doc
(
    id         text                                      not null
        primary key,
    project_id text                     default ''::text not null,
    folder_id  text                     default ''::text not null,
    name       text                     default ''::text not null,
    content    text                     default ''::text not null,
    created_by text                     default ''::text not null,
    updated_by text                     default ''::text not null,
    created_at timestamp with time zone default CURRENT_TIMESTAMP,
    updated_at timestamp with time zone default CURRENT_TIMESTAMP,
    deleted_at timestamp with time zone
);

comment on column mv_doc.name is '文档标题';

comment on column mv_doc.content is '文档内容';

comment on column mv_doc.created_by is '创建人';

comment on column mv_doc.updated_by is '更新人';

alter table mv_doc
    owner to postgres;

create index idx_mv_doc_project
    on mv_doc (project_id);

create index idx_mv_doc_project_folder
    on mv_doc (project_id, folder_id);

create index idx_mv_doc_created_at
    on mv_doc (created_at);

create index idx_mv_doc_deleted_at
    on mv_doc (deleted_at);

create table website_visit
(
    id         text not null
        primary key,
    ip         text                     default ''::text,
    user_agent text                     default ''::text,
    referer    text                     default ''::text,
    page       text                     default ''::text,
    method     text                     default ''::text,
    session_id text                     default ''::text,
    user_id    text                     default ''::text,
    created_at timestamp with time zone default CURRENT_TIMESTAMP,
    updated_at timestamp with time zone default CURRENT_TIMESTAMP,
    deleted_at timestamp with time zone
);

comment on table website_visit is '官网访问记录';

comment on column website_visit.id is '主键ID';

comment on column website_visit.ip is '访问者IP地址';

comment on column website_visit.user_agent is '用户代理信息';

comment on column website_visit.referer is '来源页面';

comment on column website_visit.page is '访问页面路径';

comment on column website_visit.method is 'HTTP请求方法';

comment on column website_visit.session_id is '会话ID';

comment on column website_visit.user_id is '用户ID（如果已登录）';

comment on column website_visit.created_at is '创建时间';

comment on column website_visit.updated_at is '更新时间';

comment on column website_visit.deleted_at is '删除时间';

alter table website_visit
    owner to postgres;

