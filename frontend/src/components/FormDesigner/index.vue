<template>
  <div class="form-designer">
    <!-- 顶部工具栏 -->
    <div class="designer-header">
      <div class="header-left">
        <t-input v-model="formConfig.title" placeholder="填写表单" class="form-title-input" />
        <t-tag theme="default" variant="light">表单</t-tag>
      </div>
      <div class="header-right">
        <!-- 尺寸切换 -->
        <t-radio-group v-model="canvasSize" variant="default-filled" size="small" class="size-switch">
          <t-radio-button value="pc">
            <t-icon name="desktop" />
            PC
          </t-radio-button>
          <t-radio-button value="mobile">
            <t-icon name="mobile" />
            移动
          </t-radio-button>
        </t-radio-group>
        
        <t-divider layout="vertical" />
        
        <t-button variant="text" @click="showStats = true">
          <t-icon name="chart-bar" />
          统计
        </t-button>
        <t-button variant="text" @click="showTheme = true">
          <t-icon name="palette" />
          主题
        </t-button>
        <t-button variant="text" @click="showSettings = true">
          <t-icon name="setting" />
          设置
        </t-button>
        <t-dropdown trigger="click" :options="shareDropdownOptions" @click="handleShareAction">
          <t-button theme="primary" class="share-btn">
            <template #icon>
              <t-icon name="share" />
            </template>
            分享
          </t-button>
        </t-dropdown>
      </div>
    </div>

    <!-- 主体区域 -->
    <div class="designer-body">
      <!-- 左侧字段库 -->
      <div class="field-library">
        <div class="library-header">
          <span class="header-title">可选题目</span>
          <div class="header-actions">
            <t-tooltip content="全部添加">
              <t-button size="small" variant="text" @click="handleAddAllFields">
                <t-icon name="add-circle" />
              </t-button>
            </t-tooltip>
            <t-tooltip content="全部隐藏">
              <t-button size="small" variant="text" @click="handleHideAllFields">
                <t-icon name="remove" />
              </t-button>
            </t-tooltip>
          </div>
        </div>

        <div class="library-content">
          <!-- 基础题型 -->
          <div class="field-group">
            <div class="group-header" @click="toggleGroup('basic')">
              <t-icon 
                :name="expandedGroups.basic ? 'chevron-down' : 'chevron-right'" 
                class="expand-icon" 
                size="16px"
              />
              <span class="group-title">基础题型</span>
              <span class="field-count">{{ basicFields.length }}</span>
            </div>
            <t-collapse-transition>
              <div v-show="expandedGroups.basic" class="field-list">
                <div
                  v-for="field in basicFields"
                  :key="field.type"
                  class="field-item"
                  draggable="true"
                  @dragstart="handleDragStart(field, $event)"
                  @dragend="handleDragEnd"
                >
                  <div class="field-icon">
                    <!-- 使用TDesign确定存在的图标 -->
                    <template v-if="field.type === 'text'">
                      <svg width="18" height="18" viewBox="0 0 24 24" fill="currentColor">
                        <path d="M3 17.25V21h3.75L17.81 9.94l-3.75-3.75L3 17.25zM20.71 7.04c.39-.39.39-1.02 0-1.41l-2.34-2.34c-.39-.39-1.02-.39-1.41 0l-1.83 1.83 3.75 3.75 1.83-1.83z"/>
                      </svg>
                    </template>
                    <template v-else-if="field.type === 'richText'">
                      <svg width="18" height="18" viewBox="0 0 24 24" fill="currentColor">
                        <path d="M14 2H6c-1.1 0-1.99.9-1.99 2L4 20c0 1.1.89 2 1.99 2H18c1.1 0 2-.9 2-2V8l-6-6zm2 16H8v-2h8v2zm0-4H8v-2h8v2zm-3-5V3.5L18.5 9H13z"/>
                      </svg>
                    </template>
                    <template v-else-if="field.type === 'singleChoice'">
                      <svg width="18" height="18" viewBox="0 0 24 24" fill="currentColor">
                        <path d="M12 2C6.48 2 2 6.48 2 12s4.48 10 10 10 10-4.48 10-10S17.52 2 12 2zm-2 15l-5-5 1.41-1.41L10 14.17l7.59-7.59L19 8l-9 9z"/>
                      </svg>
                    </template>
                    <template v-else-if="field.type === 'multipleChoice'">
                      <svg width="18" height="18" viewBox="0 0 24 24" fill="currentColor">
                        <path d="M19 3H5c-1.11 0-2 .9-2 2v14c0 1.1.89 2 2 2h14c1.11 0 2-.9 2-2V5c0-1.1-.89-2-2-2zm-9 14l-5-5 1.41-1.41L10 14.17l7.59-7.59L19 8l-9 9z"/>
                      </svg>
                    </template>
                    <template v-else-if="field.type === 'file'">
                      <svg width="18" height="18" viewBox="0 0 24 24" fill="currentColor">
                        <path d="M16.5 6v11.5c0 2.21-1.79 4-4 4s-4-1.79-4-4V5c0-1.38 1.12-2.5 2.5-2.5s2.5 1.12 2.5 2.5v10.5c0 .55-.45 1-1 1s-1-.45-1-1V6H10v9.5c0 1.38 1.12 2.5 2.5 2.5s2.5-1.12 2.5-2.5V5c0-2.21-1.79-4-4-4S7 2.79 7 5v12.5c0 3.04 2.46 5.5 5.5 5.5s5.5-2.46 5.5-5.5V6h-1.5z"/>
                      </svg>
                    </template>
                  </div>
                  <span class="field-label">{{ field.label }}</span>
                  <t-icon name="add" class="add-icon" size="16px" @click="handleQuickAdd(field)" />
                </div>
              </div>
            </t-collapse-transition>
          </div>

          <!-- 常用题型 -->
          <div class="field-group">
            <div class="group-header" @click="toggleGroup('common')">
              <t-icon 
                :name="expandedGroups.common ? 'chevron-down' : 'chevron-right'" 
                class="expand-icon" 
                size="16px"
              />
              <span class="group-title">常用题型</span>
              <span class="field-count">{{ commonFields.length }}</span>
            </div>
            <t-collapse-transition>
              <div v-show="expandedGroups.common" class="field-list">
                <div
                  v-for="field in commonFields"
                  :key="field.type"
                  class="field-item"
                  draggable="true"
                  @dragstart="handleDragStart(field, $event)"
                  @dragend="handleDragEnd"
                >
                  <div class="field-icon">
                    <!-- 使用TDesign确定存在的图标 -->
                    <template v-if="field.type === 'person'">
                      <svg width="18" height="18" viewBox="0 0 24 24" fill="currentColor">
                        <path d="M12 12c2.21 0 4-1.79 4-4s-1.79-4-4-4-4 1.79-4 4 1.79 4 4 4zm0 2c-2.67 0-8 1.34-8 4v2h16v-2c0-2.66-5.33-4-8-4z"/>
                      </svg>
                    </template>
                    <template v-else-if="field.type === 'department'">
                      <svg width="18" height="18" viewBox="0 0 24 24" fill="currentColor">
                        <path d="M16 11c1.66 0 2.99-1.34 2.99-3S17.66 5 16 5c-1.66 0-3 1.34-3 3s1.34 3 3 3zm-8 0c1.66 0 2.99-1.34 2.99-3S9.66 5 8 5C6.34 5 5 6.34 5 8s1.34 3 3 3zm0 2c-2.33 0-7 1.17-7 3.5V19h14v-2.5c0-2.33-4.67-3.5-7-3.5zm8 0c-.29 0-.62.02-.97.05 1.16.84 1.97 1.97 1.97 3.45V19h6v-2.5c0-2.33-4.67-3.5-7-3.5z"/>
                      </svg>
                    </template>
                    <template v-else-if="field.type === 'date'">
                      <svg width="18" height="18" viewBox="0 0 24 24" fill="currentColor">
                        <path d="M19 3h-1V1h-2v2H8V1H6v2H5c-1.11 0-1.99.9-1.99 2L3 19c0 1.1.89 2 2 2h14c1.1 0 2-.9 2-2V5c0-1.1-.9-2-2-2zm0 16H5V8h14v11zM7 10h5v5H7z"/>
                      </svg>
                    </template>
                    <template v-else-if="field.type === 'number'">
                      <svg width="18" height="18" viewBox="0 0 24 24" fill="currentColor">
                        <path d="M19 3H5c-1.1 0-2 .9-2 2v14c0 1.1.9 2 2 2h14c1.1 0 2-.9 2-2V5c0-1.1-.9-2-2-2zm-5 14h-2v-4H8v-2h4V7h2v4h4v2h-4v4z"/>
                      </svg>
                    </template>
                    <template v-else-if="field.type === 'checkbox'">
                      <svg width="18" height="18" viewBox="0 0 24 24" fill="currentColor">
                        <path d="M19 5v14H5V5h14m0-2H5c-1.1 0-2 .9-2 2v14c0 1.1.9 2 2 2h14c1.1 0 2-.9 2-2V5c0-1.1-.9-2-2-2z"/>
                      </svg>
                    </template>
                    <template v-else-if="field.type === 'link'">
                      <svg width="18" height="18" viewBox="0 0 24 24" fill="currentColor">
                        <path d="M3.9 12c0-1.71 1.39-3.1 3.1-3.1h4V7H7c-2.76 0-5 2.24-5 5s2.24 5 5 5h4v-1.9H7c-1.71 0-3.1-1.39-3.1-3.1zM8 13h8v-2H8v2zm9-6h-4v1.9h4c1.71 0 3.1 1.39 3.1 3.1s-1.39 3.1-3.1 3.1h-4V17h4c2.76 0 5-2.24 5-5s-2.24-5-5-5z"/>
                      </svg>
                    </template>
                    <template v-else-if="field.type === 'location'">
                      <svg width="18" height="18" viewBox="0 0 24 24" fill="currentColor">
                        <path d="M12 2C8.13 2 5 5.13 5 9c0 5.25 7 13 7 13s7-7.75 7-13c0-3.87-3.13-7-7-7zm0 9.5c-1.38 0-2.5-1.12-2.5-2.5s1.12-2.5 2.5-2.5 2.5 1.12 2.5 2.5-1.12 2.5-2.5 2.5z"/>
                      </svg>
                    </template>
                    <template v-else-if="field.type === 'phone'">
                      <svg width="18" height="18" viewBox="0 0 24 24" fill="currentColor">
                        <path d="M6.62 10.79c1.44 2.83 3.76 5.14 6.59 6.59l2.2-2.2c.27-.27.67-.36 1.02-.24 1.12.37 2.33.57 3.57.57.55 0 1 .45 1 1V20c0 .55-.45 1-1 1-9.39 0-17-7.61-17-17 0-.55.45-1 1-1h3.5c.55 0 1 .45 1 1 0 1.25.2 2.45.57 3.57.11.35.03.74-.25 1.02l-2.2 2.2z"/>
                      </svg>
                    </template>
                    <template v-else-if="field.type === 'area'">
                      <svg width="18" height="18" viewBox="0 0 24 24" fill="currentColor">
                        <path d="M20.5 3l-.16.03L15 5.1 9 3 3.36 4.9c-.21.07-.36.25-.36.48V20.5c0 .28.22.5.5.5l.16-.03L9 18.9l6 2.1 5.64-1.9c.21-.07.36-.25.36-.48V3.5c0-.28-.22-.5-.5-.5zM15 19l-6-2.11V5l6 2.11V19z"/>
                      </svg>
                    </template>
                    <template v-else-if="field.type === 'rating'">
                      <svg width="18" height="18" viewBox="0 0 24 24" fill="currentColor">
                        <path d="M12 17.27L18.18 21l-1.64-7.03L22 9.24l-7.19-.61L12 2 9.19 8.63 2 9.24l5.46 4.73L5.82 21z"/>
                      </svg>
                    </template>
                  </div>
                  <span class="field-label">{{ field.label }}</span>
                  <t-icon name="add" class="add-icon" size="16px" @click="handleQuickAdd(field)" />
                </div>
              </div>
            </t-collapse-transition>
          </div>
        </div>
      </div>

      <!-- 中间画布区 -->
      <div class="form-canvas">
        <div class="canvas-wrapper" :class="`canvas-${canvasSize}`">
          <!-- 表单封面 -->
          <div class="form-cover" :style="coverStyle">
            <div class="cover-content">
              <div class="cover-icon">
                <t-icon name="star-filled" size="48px" />
              </div>
            </div>
            <t-button variant="text" class="edit-cover-btn" @click="showCoverEditor = true">
              <t-icon name="image" />
            </t-button>
          </div>

          <!-- 表单标题 -->
          <div class="form-header-section">
            <div class="form-title-wrapper">
              <t-input
                v-model="formConfig.title"
                placeholder="请输入表单标题"
                class="form-title-input"
                :bordered="false"
                @focus="handleTitleFocus"
                @blur="handleTitleBlur"
              />
            </div>
            <div class="form-description-wrapper">
              <t-textarea
                v-model="formConfig.description"
                placeholder="请输入表单描述（可选）"
                class="form-description-input"
                :bordered="false"
                :autosize="{ minRows: 1, maxRows: 3 }"
                @focus="handleDescFocus"
                @blur="handleDescBlur"
              />
            </div>
          </div>

          <!-- 表单字段列表 -->
          <div
            ref="formFieldsRef"
            class="form-fields"
            @dragover.prevent="handleCanvasDragOver"
            @drop="handleCanvasDrop"
            @dragleave="handleCanvasDragLeave"
          >
            <!-- 拖拽插入指示线 -->
            <div 
              v-if="dropIndicatorIndex !== -1" 
              class="drop-indicator"
              :style="{ top: dropIndicatorTop + 'px' }"
            ></div>

            <div
              v-for="(field, index) in formFields"
              :key="field.id"
              class="form-field-item"
              :class="{ 
                active: selectedFieldId === field.id,
                'is-dragging': draggingFieldId === field.id
              }"
              :data-index="index"
              draggable="true"
              @click="handleSelectField(field)"
              @dragstart="handleFieldDragStart(field, index, $event)"
              @dragover.prevent.stop="handleFieldItemDragOver(index, $event)"
              @dragleave.stop="handleFieldDragLeave"
              @drop.prevent.stop="handleFieldItemDrop(index, $event)"
              @dragend="handleFieldDragEnd"
            >
              <div class="field-drag-handle" @mousedown.stop>
                <div class="drag-handle-icon">
                  <span></span>
                  <span></span>
                  <span></span>
                  <span></span>
                  <span></span>
                  <span></span>
                </div>
              </div>
              <div class="field-content">
                <div class="field-header">
                  <span class="field-label">
                    <span v-if="field.required" class="required-mark">*</span>
                    {{ field.config.label || field.label }}
                  </span>
                  <div class="field-actions">
                    <t-tooltip content="复制字段">
                      <t-button size="small" variant="text" @click.stop="handleCopyField(index)">
                        <t-icon name="file-copy" />
                      </t-button>
                    </t-tooltip>
                    <t-tooltip content="删除字段">
                      <t-button size="small" variant="text" @click.stop="handleDeleteField(index)">
                        <t-icon name="delete" />
                      </t-button>
                    </t-tooltip>
                  </div>
                </div>
                <div class="field-input">
                  <FieldPreview
                    :field="field"
                    :readonly="true"
                  />
                </div>
              </div>
            </div>

            <!-- 空状态 -->
            <div v-if="formFields.length === 0" class="empty-canvas">
              <div class="empty-icon">
                <t-icon name="view-module" size="64px" />
              </div>
              <p class="empty-title">开始设计你的表单</p>
              <p class="empty-desc">从左侧拖拽字段到此处，或点击字段快速添加</p>
            </div>
          </div>
        </div>
      </div>

      <!-- 右侧配置区 -->
      <div class="field-config-panel" v-if="selectedField">
        <div class="panel-header">
          <span>字段设置</span>
          <t-button size="small" variant="text" @click="selectedFieldId = null">
            <t-icon name="close" />
          </t-button>
        </div>

        <div class="panel-content">
          <!-- 字段标签 -->
          <div class="config-item">
            <div class="config-label">字段标题</div>
            <t-input v-model="selectedField.config.label" placeholder="请输入字段标题" />
          </div>

          <!-- 字段描述 -->
          <div class="config-item">
            <div class="config-label">字段描述</div>
            <t-textarea
              v-model="selectedField.config.description"
              placeholder="请输入字段描述"
              :autosize="{ minRows: 2, maxRows: 4 }"
            />
          </div>

          <!-- 必填选项 -->
          <div class="config-item">
            <t-checkbox v-model="selectedField.required">必填</t-checkbox>
          </div>

          <!-- 隐藏选项 -->
          <div class="config-item">
            <t-checkbox v-model="selectedField.config.hidden">隐藏</t-checkbox>
          </div>

          <!-- 默认值 -->
          <div class="config-item" v-if="!['file', 'image'].includes(selectedField.type)">
            <div class="config-label">默认值</div>
            <t-input v-model="selectedField.config.defaultValue" placeholder="请输入默认值" />
          </div>

          <!-- 选项配置（针对单选、多选） -->
          <div class="config-item" v-if="['singleChoice', 'multipleChoice'].includes(selectedField.type)">
            <div class="config-label">选项配置</div>
            <div class="options-list">
              <div v-for="(option, idx) in selectedField.config.options" :key="idx" class="option-item">
                <t-input v-model="option.label" placeholder="选项内容" size="small" />
                <t-button size="small" variant="text" @click="handleDeleteOption(idx)">
                  <t-icon name="delete" />
                </t-button>
              </div>
              <t-button size="small" variant="dashed" block @click="handleAddOption">
                <t-icon name="add" />
                添加选项
              </t-button>
            </div>
          </div>

          <!-- 数字范围（针对数字类型） -->
          <div class="config-item" v-if="selectedField.type === 'number'">
            <div class="config-label">数值范围</div>
            <div class="range-inputs">
              <t-input-number v-model="selectedField.config.min" placeholder="最小值" size="small" />
              <span>-</span>
              <t-input-number v-model="selectedField.config.max" placeholder="最大值" size="small" />
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- 统计弹窗 -->
    <t-dialog
      v-model:visible="showStats"
      header="表单统计"
      width="800px"
      :footer="false"
    >
      <div class="stats-content">
        <div class="stats-overview">
          <div class="stat-card">
            <div class="stat-icon" style="background: #e7f5ff;">
              <t-icon name="usergroup" size="24px" style="color: #1890ff;" />
            </div>
            <div class="stat-info">
              <div class="stat-value">{{ statsData.totalSubmissions }}</div>
              <div class="stat-label">总提交数</div>
            </div>
          </div>
          <div class="stat-card">
            <div class="stat-icon" style="background: #fff1f0;">
              <t-icon name="chart-line" size="24px" style="color: #f5222d;" />
            </div>
            <div class="stat-info">
              <div class="stat-value">{{ statsData.todaySubmissions }}</div>
              <div class="stat-label">今日提交</div>
            </div>
          </div>
          <div class="stat-card">
            <div class="stat-icon" style="background: #f6ffed;">
              <t-icon name="chart-pie" size="24px" style="color: #52c41a;" />
            </div>
            <div class="stat-info">
              <div class="stat-value">{{ statsData.completionRate }}%</div>
              <div class="stat-label">完成率</div>
            </div>
          </div>
          <div class="stat-card">
            <div class="stat-icon" style="background: #fff7e6;">
              <t-icon name="time" size="24px" style="color: #faad14;" />
            </div>
            <div class="stat-info">
              <div class="stat-value">{{ statsData.avgTime }}s</div>
              <div class="stat-label">平均时长</div>
            </div>
          </div>
        </div>

        <div class="stats-chart">
          <div class="chart-header">
            <span class="chart-title">提交趋势</span>
            <t-radio-group v-model="chartTimeRange" variant="default-filled" size="small">
              <t-radio-button value="7days">7天</t-radio-button>
              <t-radio-button value="30days">30天</t-radio-button>
              <t-radio-button value="90days">90天</t-radio-button>
            </t-radio-group>
          </div>
          <div class="chart-placeholder">
            <t-icon name="chart-bar" size="48px" />
            <p>图表区域（可集成ECharts）</p>
          </div>
        </div>

        <div class="stats-actions">
          <t-button theme="primary" @click="handleExportStats">
            <t-icon name="download" />
            导出数据
          </t-button>
        </div>
      </div>
    </t-dialog>

    <!-- 主题设置弹窗 -->
    <t-dialog
      v-model:visible="showTheme"
      header="主题设置"
      width="680px"
      @confirm="handleSaveTheme"
    >
      <div class="theme-content">
        <div class="theme-section">
          <div class="section-title">
            <span class="title-text">封面主题</span>
            <span class="title-desc">选择你喜欢的封面颜色</span>
          </div>
          <div class="gradient-grid">
            <div
              v-for="(gradient, idx) in gradientPresets"
              :key="idx"
              class="gradient-card"
              :class="{ selected: formConfig.coverGradient === gradient }"
              @click="formConfig.coverGradient = gradient"
            >
              <div class="gradient-preview" :style="{ background: gradient }">
                <div v-if="formConfig.coverGradient === gradient" class="selected-badge">
                  <t-icon name="check" size="16px" />
                </div>
              </div>
            </div>
          </div>
        </div>

        <div class="theme-section">
          <div class="section-title">
            <span class="title-text">主题颜色</span>
            <span class="title-desc">设置表单的主色调</span>
          </div>
          <div class="color-grid">
            <div
              v-for="color in themeColors"
              :key="color.value"
              class="color-card"
              :class="{ selected: formConfig.themeColor === color.value }"
              @click="selectThemeColor(color.value)"
            >
              <div class="color-preview" :style="{ background: color.value }">
                <div v-if="formConfig.themeColor === color.value" class="selected-badge">
                  <t-icon name="check" size="16px" />
                </div>
              </div>
              <div class="color-name">{{ color.name }}</div>
            </div>
          </div>
        </div>
      </div>
    </t-dialog>

    <!-- 设置弹窗 -->
    <t-dialog
      v-model:visible="showSettings"
      header="表单设置"
      width="600px"
      @confirm="handleSaveSettings"
    >
      <div class="settings-content">
        <div class="setting-item">
          <div class="setting-label">
            <span>允许多次提交</span>
            <span class="setting-desc">用户可以多次填写并提交此表单</span>
          </div>
          <t-switch v-model="formSettings.allowMultipleSubmit" />
        </div>

        <div class="setting-item">
          <div class="setting-label">
            <span>提交后显示结果</span>
            <span class="setting-desc">用户提交后可以查看提交结果</span>
          </div>
          <t-switch v-model="formSettings.showResultAfterSubmit" />
        </div>

        <div class="setting-item">
          <div class="setting-label">
            <span>需要登录</span>
            <span class="setting-desc">用户需要登录后才能填写表单</span>
          </div>
          <t-switch v-model="formSettings.requireLogin" />
        </div>

        <div class="setting-item">
          <div class="setting-label">
            <span>收集微信信息</span>
            <span class="setting-desc">在微信中打开时自动获取用户信息</span>
          </div>
          <t-switch v-model="formSettings.collectWechatInfo" />
        </div>

        <t-divider />

        <div class="setting-item">
          <div class="setting-label">
            <span>截止时间</span>
            <span class="setting-desc">设置表单的截止时间</span>
          </div>
          <t-date-picker
            v-model="formSettings.deadline"
            enable-time-picker
            placeholder="选择截止时间"
          />
        </div>

        <div class="setting-item">
          <div class="setting-label">
            <span>提交限制</span>
            <span class="setting-desc">限制最大提交数量</span>
          </div>
          <t-input-number
            v-model="formSettings.maxSubmissions"
            :min="0"
            placeholder="0表示不限制"
          />
        </div>
      </div>
    </t-dialog>

    <!-- 分享弹窗 -->
    <t-dialog
      v-model:visible="showShareDialog"
      header="分享表单"
      width="600px"
      :footer="false"
    >
      <div class="share-content">
        <t-tabs v-model="shareTab">
          <t-tab-panel value="link" label="链接分享">
            <div class="share-link-section">
              <div class="link-display">
                <t-input
                  v-model="shareLink"
                  readonly
                  placeholder="表单链接"
                >
                  <template #suffix>
                    <t-button size="small" variant="text" @click="handleCopyLink">
                      <t-icon name="file-copy" />
                      复制
                    </t-button>
                  </template>
                </t-input>
              </div>

              <div class="qrcode-section">
                <div class="qrcode-title">扫码填写</div>
                <div class="qrcode-placeholder">
                  <t-icon name="qrcode" size="80px" />
                  <p>二维码区域</p>
                </div>
              </div>

              <div class="share-options">
                <t-checkbox v-model="shareOptionsData.allowEdit">允许编辑</t-checkbox>
                <t-checkbox v-model="shareOptionsData.allowView">允许查看结果</t-checkbox>
              </div>
            </div>
          </t-tab-panel>

          <t-tab-panel value="embed" label="嵌入代码">
            <div class="embed-section">
              <t-textarea
                v-model="embedCode"
                readonly
                :autosize="{ minRows: 6, maxRows: 10 }"
                placeholder="嵌入代码"
              />
              <t-button block theme="primary" @click="handleCopyEmbedCode">
                <t-icon name="file-copy" />
                复制代码
              </t-button>
            </div>
          </t-tab-panel>
        </t-tabs>
      </div>
    </t-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { MessagePlugin } from 'tdesign-vue-next'
import FieldPreview from './FieldPreview.vue'

interface FieldConfig {
  label?: string
  description?: string
  defaultValue?: any
  hidden?: boolean
  options?: Array<{ label: string; value: string }>
  min?: number
  max?: number
}

interface FormField {
  id: string
  type: string
  label: string
  icon: string
  required: boolean
  config: FieldConfig
}

const props = defineProps<{
  tableId: string
  fields: any[]
}>()

const emit = defineEmits(['save'])

// 表单配置
const formConfig = ref({
  title: '填写表单',
  description: '',
  coverGradient: 'linear-gradient(135deg, #667eea 0%, #764ba2 100%)',
  themeColor: '#0052d9',
})

// 表单字段列表
const formFields = ref<FormField[]>([])

// 选中的字段ID
const selectedFieldId = ref<string | null>(null)

// 拖拽相关
const dragOverIndex = ref<number>(-1)
const dragFieldIndex = ref<number>(-1)
const draggingFieldId = ref<string | null>(null)
const dropIndicatorIndex = ref<number>(-1)
const dropIndicatorTop = ref<number>(0)
const formFieldsRef = ref<HTMLElement | null>(null)

// 画布尺寸
const canvasSize = ref<'pc' | 'mobile'>('pc')

// 弹窗状态
const showStats = ref(false)
const showTheme = ref(false)
const showSettings = ref(false)
const showCoverEditor = ref(false)
const showShareDialog = ref(false)

// 分享相关
const shareTab = ref('link')
const shareLink = ref('https://form.example.com/f/abc123')
const embedCode = ref('<iframe src="https://form.example.com/f/abc123" width="100%" height="600"></iframe>')

// 分享下拉选项
const shareDropdownOptions = ref([
  { content: '复制链接', value: 'copy' },
  { content: '生成二维码', value: 'qrcode' },
  { content: '嵌入网页', value: 'embed' },
])

// 分享选项
const shareOptionsData = ref({
  allowEdit: false,
  allowView: true,
})

// 字段组展开状态
const expandedGroups = ref({
  basic: true,
  common: true,
})

// 图表时间范围
const chartTimeRange = ref('7days')

// 统计数据
const statsData = ref({
  totalSubmissions: 128,
  todaySubmissions: 15,
  completionRate: 87,
  avgTime: 125,
})

// 表单设置
const formSettings = ref({
  allowMultipleSubmit: true,
  showResultAfterSubmit: false,
  requireLogin: false,
  collectWechatInfo: false,
  deadline: '',
  maxSubmissions: 0,
})

// 基础字段类型
const basicFields = [
  { type: 'text', label: '文本', icon: 'edit' },
  { type: 'richText', label: '富文本', icon: 'file-paste' },
  { type: 'singleChoice', label: '单选', icon: 'check-circle' },
  { type: 'multipleChoice', label: '多选', icon: 'check' },
  { type: 'file', label: '图片和附件', icon: 'attach' },
]

// 常用字段类型
const commonFields = [
  { type: 'person', label: '人员', icon: 'user' },
  { type: 'department', label: '部门', icon: 'usergroup' },
  { type: 'date', label: '日期', icon: 'calendar' },
  { type: 'number', label: '数字', icon: 'sort-descending' },
  { type: 'checkbox', label: '复选框', icon: 'check-rectangle' },
  { type: 'link', label: '链接', icon: 'link' },
  { type: 'location', label: '地理位置', icon: 'location' },
  { type: 'phone', label: '电话', icon: 'call' },
  { type: 'area', label: '行政区域', icon: 'map-location' },
  { type: 'rating', label: '评分', icon: 'star' },
]

// 主题颜色选项
const themeColors = [
  { name: '蓝色', value: '#0052d9' },
  { name: '绿色', value: '#00a870' },
  { name: '橙色', value: '#e37318' },
  { name: '红色', value: '#c9353f' },
  { name: '紫色', value: '#8b5cf6' },
  { name: '青色', value: '#06b6d4' },
]

// 渐变背景预设
const gradientPresets = [
  'linear-gradient(135deg, #667eea 0%, #764ba2 100%)',
  'linear-gradient(135deg, #f093fb 0%, #f5576c 100%)',
  'linear-gradient(135deg, #4facfe 0%, #00f2fe 100%)',
  'linear-gradient(135deg, #43e97b 0%, #38f9d7 100%)',
  'linear-gradient(135deg, #fa709a 0%, #fee140 100%)',
  'linear-gradient(135deg, #30cfd0 0%, #330867 100%)',
]

// 封面样式
const coverStyle = computed(() => ({
  background: formConfig.value.coverGradient,
}))

// 选中的字段
const selectedField = computed(() => {
  return formFields.value.find((f) => f.id === selectedFieldId.value)
})

// 选择主题颜色
const selectThemeColor = (color: string) => {
  formConfig.value.themeColor = color
  // 应用主题色到画布
  applyThemeColor(color)
}

// 应用主题颜色
const applyThemeColor = (color: string) => {
  // 更新CSS变量
  document.documentElement.style.setProperty('--form-theme-color', color)
}

// 切换字段组展开状态
const toggleGroup = (group: 'basic' | 'common') => {
  expandedGroups.value[group] = !expandedGroups.value[group]
}

// 标题焦点事件
const handleTitleFocus = () => {
  // 标题获得焦点
}

const handleTitleBlur = () => {
  // 标题失去焦点
}

// 描述焦点事件
const handleDescFocus = () => {
  // 描述获得焦点
}

const handleDescBlur = () => {
  // 描述失去焦点
}

// 快速添加字段
const handleQuickAdd = (fieldData: any) => {
  const newField: FormField = {
    id: `field_${Date.now()}`,
    type: fieldData.type,
    label: fieldData.label,
    icon: fieldData.icon,
    required: false,
    config: {
      label: fieldData.label,
      description: '',
      defaultValue: '',
      hidden: false,
      options: fieldData.type.includes('Choice') ? [
        { label: '选项一', value: '1' },
        { label: '选项二', value: '2' },
      ] : undefined,
    },
  }
  
  formFields.value.push(newField)
  MessagePlugin.success(`已添加${fieldData.label}字段`)
}

// 拖拽开始
const handleDragStart = (field: any, event: DragEvent) => {
  event.dataTransfer!.effectAllowed = 'copy'
  event.dataTransfer!.setData('field', JSON.stringify(field))
  event.dataTransfer!.setData('source', 'library')
  
  // 设置拖拽图片 - 使用画布宽度
  const dragImage = document.createElement('div')
  dragImage.className = 'drag-ghost'
  const canvasWidth = canvasSize.value === 'mobile' ? 375 : 800
  dragImage.innerHTML = `
    <div style="
      width: ${canvasWidth - 64}px;
      padding: 16px;
      background: white;
      border: 2px dashed #0052d9;
      border-radius: 8px;
      box-shadow: 0 4px 12px rgba(0, 82, 217, 0.3);
      font-size: 14px;
      color: #333;
      white-space: nowrap;
      overflow: hidden;
      text-overflow: ellipsis;
    ">
      <span style="margin-right: 8px;">📋</span>
      ${field.label}
    </div>
  `
  document.body.appendChild(dragImage)
  event.dataTransfer!.setDragImage(dragImage, 0, 0)
  setTimeout(() => document.body.removeChild(dragImage), 0)
}

// 拖拽结束
const handleDragEnd = () => {
  // 清理拖拽状态
}

// 字段拖拽开始
const handleFieldDragStart = (field: FormField, index: number, event: DragEvent) => {
  dragFieldIndex.value = index
  draggingFieldId.value = field.id
  event.dataTransfer!.effectAllowed = 'move'
  event.dataTransfer!.setData('fieldIndex', String(index))
  event.dataTransfer!.setData('source', 'canvas')
  
  // 设置拖拽图片
  const dragImage = document.createElement('div')
  dragImage.className = 'drag-ghost'
  dragImage.innerHTML = `
    <div style="
      padding: 12px 16px;
      background: white;
      border: 2px solid #0052d9;
      border-radius: 8px;
      box-shadow: 0 4px 12px rgba(0, 82, 217, 0.3);
      font-size: 14px;
      color: #333;
      white-space: nowrap;
    ">
      <span style="margin-right: 8px;">🔄</span>
      ${field.config.label || field.label}
    </div>
  `
  document.body.appendChild(dragImage)
  event.dataTransfer!.setDragImage(dragImage, 0, 0)
  setTimeout(() => document.body.removeChild(dragImage), 0)
}

// 画布拖拽悬停
const handleCanvasDragOver = (event: DragEvent) => {
  event.preventDefault()
  const source = event.dataTransfer!.types.includes('source')
  if (!source) {
    updateDropIndicator(event)
  }
}

// 画布拖拽离开
const handleCanvasDragLeave = (event: DragEvent) => {
  const relatedTarget = event.relatedTarget as HTMLElement
  if (!relatedTarget || !formFieldsRef.value?.contains(relatedTarget)) {
    dropIndicatorIndex.value = -1
  }
}

// 更新拖拽插入指示线位置
const updateDropIndicator = (event: DragEvent) => {
  if (!formFieldsRef.value) return
  
  const fields = Array.from(formFieldsRef.value.querySelectorAll('.form-field-item'))
  if (fields.length === 0) {
    dropIndicatorIndex.value = 0
    dropIndicatorTop.value = 0
    return
  }
  
  const mouseY = event.clientY
  const containerRect = formFieldsRef.value.getBoundingClientRect()
  
  let insertIndex = fields.length
  let insertTop = 0
  
  for (let i = 0; i < fields.length; i++) {
    const field = fields[i] as HTMLElement
    const rect = field.getBoundingClientRect()
    const midPoint = rect.top + rect.height / 2
    
    if (mouseY < midPoint) {
      insertIndex = i
      insertTop = rect.top - containerRect.top
      break
    } else if (i === fields.length - 1) {
      insertIndex = fields.length
      insertTop = rect.bottom - containerRect.top
    }
  }
  
  dropIndicatorIndex.value = insertIndex
  dropIndicatorTop.value = insertTop
}

// 字段项拖拽悬停
const handleFieldItemDragOver = (index: number, event: DragEvent) => {
  event.preventDefault()
  event.stopPropagation()
  updateDropIndicator(event)
}

// 字段拖拽离开
const handleFieldDragLeave = () => {
  // 保留指示线，只在离开整个画布时清除
}

// 字段项拖拽放置
const handleFieldItemDrop = (targetIndex: number, event: DragEvent) => {
  event.preventDefault()
  event.stopPropagation()
  
  const source = event.dataTransfer!.getData('source')
  
  if (source === 'canvas') {
    // 画布内拖拽排序
    const fromIndex = dragFieldIndex.value
    if (fromIndex !== -1) {
      const actualTargetIndex = dropIndicatorIndex.value
      if (fromIndex !== actualTargetIndex) {
        const field = formFields.value[fromIndex]
        formFields.value.splice(fromIndex, 1)
        const newIndex = fromIndex < actualTargetIndex ? actualTargetIndex - 1 : actualTargetIndex
        formFields.value.splice(newIndex, 0, field)
        MessagePlugin.success('字段位置已调整')
      }
    }
  } else if (source === 'library') {
    // 从左侧库拖拽到画布
    const fieldData = JSON.parse(event.dataTransfer!.getData('field'))
    const newField: FormField = {
      id: `field_${Date.now()}`,
      type: fieldData.type,
      label: fieldData.label,
      icon: fieldData.icon,
      required: false,
      config: {
        label: fieldData.label,
        description: '',
        defaultValue: '',
        hidden: false,
        options: fieldData.type.includes('Choice') ? [
          { label: '选项一', value: '1' },
          { label: '选项二', value: '2' },
        ] : undefined,
      },
    }
    
    const insertIndex = dropIndicatorIndex.value
    if (insertIndex >= 0 && insertIndex <= formFields.value.length) {
      formFields.value.splice(insertIndex, 0, newField)
    } else {
      formFields.value.push(newField)
    }
    MessagePlugin.success(`已添加${fieldData.label}字段`)
  }
  
  dropIndicatorIndex.value = -1
  dragOverIndex.value = -1
  dragFieldIndex.value = -1
  draggingFieldId.value = null
}

// 画布拖拽放置
const handleCanvasDrop = (event: DragEvent) => {
  event.preventDefault()
  
  const source = event.dataTransfer!.getData('source')
  
  if (source === 'library') {
    const fieldData = JSON.parse(event.dataTransfer!.getData('field'))
    const newField: FormField = {
      id: `field_${Date.now()}`,
      type: fieldData.type,
      label: fieldData.label,
      icon: fieldData.icon,
      required: false,
      config: {
        label: fieldData.label,
        description: '',
        defaultValue: '',
        hidden: false,
        options: fieldData.type.includes('Choice') ? [
          { label: '选项一', value: '1' },
          { label: '选项二', value: '2' },
        ] : undefined,
      },
    }
    
    const insertIndex = dropIndicatorIndex.value
    if (insertIndex >= 0 && insertIndex <= formFields.value.length) {
      formFields.value.splice(insertIndex, 0, newField)
    } else {
      formFields.value.push(newField)
    }
    MessagePlugin.success(`已添加${fieldData.label}字段`)
  }
  
  dropIndicatorIndex.value = -1
}

// 字段拖拽结束
const handleFieldDragEnd = () => {
  dragOverIndex.value = -1
  dragFieldIndex.value = -1
  draggingFieldId.value = null
  dropIndicatorIndex.value = -1
}

// 选中字段
const handleSelectField = (field: FormField) => {
  selectedFieldId.value = field.id
}

// 复制字段
const handleCopyField = (index: number) => {
  const field = formFields.value[index]
  const newField: FormField = {
    ...field,
    id: `field_${Date.now()}`,
    config: { ...field.config },
  }
  formFields.value.splice(index + 1, 0, newField)
  MessagePlugin.success('字段已复制')
}

// 删除字段
const handleDeleteField = (index: number) => {
  formFields.value.splice(index, 1)
  selectedFieldId.value = null
  MessagePlugin.success('字段已删除')
}

// 添加所有字段
const handleAddAllFields = () => {
  MessagePlugin.info('功能开发中')
}

// 隐藏所有字段
const handleHideAllFields = () => {
  MessagePlugin.info('功能开发中')
}

const handleAddFieldGroup = (type: string) => {
  MessagePlugin.info('功能开发中')
}

const handleAddCustomField = () => {
  MessagePlugin.info('功能开发中')
}

// 添加选项
const handleAddOption = () => {
  if (selectedField.value?.config.options) {
    selectedField.value.config.options.push({
      label: `选项${selectedField.value.config.options.length + 1}`,
      value: `${selectedField.value.config.options.length + 1}`,
    })
  }
}

// 删除选项
const handleDeleteOption = (index: number) => {
  if (selectedField.value?.config.options) {
    selectedField.value.config.options.splice(index, 1)
  }
}

// 分享操作
const handleShareAction = (data: any) => {
  const action = data.value
  if (action === 'copy') {
    handleCopyLink()
  } else if (action === 'qrcode' || action === 'embed') {
    shareTab.value = action === 'qrcode' ? 'link' : 'embed'
    showShareDialog.value = true
  }
}

// 复制链接
const handleCopyLink = () => {
  navigator.clipboard.writeText(shareLink.value)
  MessagePlugin.success('链接已复制到剪贴板')
}

// 复制嵌入代码
const handleCopyEmbedCode = () => {
  navigator.clipboard.writeText(embedCode.value)
  MessagePlugin.success('代码已复制到剪贴板')
}

// 导出统计数据
const handleExportStats = () => {
  MessagePlugin.success('正在导出数据...')
  // TODO: 实现导出逻辑
}

// 保存主题
const handleSaveTheme = () => {
  applyThemeColor(formConfig.value.themeColor)
  showTheme.value = false
  MessagePlugin.success('主题已保存')
}

// 保存设置
const handleSaveSettings = () => {
  showSettings.value = false
  MessagePlugin.success('设置已保存')
}

const handleSaveCover = () => {
  showCoverEditor.value = false
  MessagePlugin.success('封面已保存')
}

// 获取字段预览组件
const getFieldPreviewComponent = () => {
  return FieldPreview
}
</script>

<style scoped lang="less">
:root {
  --form-theme-color: #0052d9;
}

.form-designer {
  height: 100vh;
  display: flex;
  flex-direction: column;
  background: #f5f7fa;

  .designer-header {
    display: flex;
    align-items: center;
    justify-content: space-between;
    padding: 12px 24px;
    background: #fff;
    border-bottom: 1px solid #e8e8e8;

    .header-left {
      display: flex;
      align-items: center;
      gap: 12px;

      .form-title-input {
        width: 200px;
      }
    }

    .header-right {
      display: flex;
      align-items: center;
      gap: 8px;
      
      .size-switch {
        margin-right: 8px;
        
        :deep(.t-radio-button) {
          display: flex;
          align-items: center;
          gap: 4px;
        }
      }
      
      .share-btn {
        :deep(.t-button__icon) {
          margin-right: 4px;
        }
      }
    }
  }

  .designer-body {
    flex: 1;
    display: flex;
    overflow: hidden;

    // 左侧字段库
    .field-library {
      width: 280px;
      background: #fff;
      border-right: 1px solid #e8e8e8;
      overflow-y: auto;
      flex-shrink: 0;
      display: flex;
      flex-direction: column;

      .library-header {
        padding: 16px;
        border-bottom: 1px solid #e8e8e8;
        display: flex;
        align-items: center;
        justify-content: space-between;
        position: sticky;
        top: 0;
        background: #fff;
        z-index: 10;

        .header-title {
          font-size: 14px;
          font-weight: 600;
          color: #333;
        }

        .header-actions {
          display: flex;
          gap: 4px;
        }
      }

      .library-content {
        flex: 1;
        overflow-y: auto;
      }

      .field-group {
        border-bottom: 1px solid #f0f0f0;

        .group-header {
          display: flex;
          align-items: center;
          gap: 8px;
          padding: 12px 16px;
          background: #fafafa;
          cursor: pointer;
          user-select: none;
          transition: all 0.2s;

          &:hover {
            background: #f5f5f5;
          }

          .expand-icon {
            color: #646a73;
            transition: all 0.2s;
            flex-shrink: 0;
          }

          .group-title {
            flex: 1;
            font-size: 13px;
            font-weight: 500;
            color: #333;
          }

          .field-count {
            font-size: 12px;
            color: #999;
            background: #f0f0f0;
            padding: 2px 8px;
            border-radius: 10px;
          }
        }

        .field-list {
          padding: 8px;

          .field-item {
            display: flex;
            align-items: center;
            gap: 12px;
            padding: 10px 12px;
            margin-bottom: 4px;
            border-radius: 6px;
            cursor: move;
            transition: all 0.2s;
            position: relative;

            &:hover {
              background: #f5f9ff;
              border: 1px solid #d9e8ff;
              padding: 9px 11px;

              .add-icon {
                opacity: 1;
              }
            }

            .field-icon {
              width: 24px;
              height: 24px;
              display: flex;
              align-items: center;
              justify-content: center;
              background: #f5f5f5;
              border-radius: 4px;
              flex-shrink: 0;
              color: #1f2329;
            }

            .field-label {
              flex: 1;
              font-size: 14px;
              color: #333;
            }

            .add-icon {
              opacity: 0;
              cursor: pointer;
              color: #0052d9;
              transition: opacity 0.2s;

              &:hover {
                color: #0034a8;
              }
            }
          }
        }
      }
    }

    // 中间画布区
    .form-canvas {
      flex: 1;
      overflow-y: auto;
      padding: 24px;
      background: #f5f7fa;
      margin-bottom: 100px;
      
      // 确保可以滚动
      &::-webkit-scrollbar {
        width: 6px;
      }
      
      &::-webkit-scrollbar-thumb {
        background-color: rgba(0, 0, 0, 0.2);
        border-radius: 3px;
        
        &:hover {
          background-color: rgba(0, 0, 0, 0.3);
        }
      }

      .canvas-wrapper {
        max-width: 800px;
        margin: 0 auto;
        background: #fff;
        border-radius: 8px;
        box-shadow: 0 2px 12px rgba(0, 0, 0, 0.08);
        overflow: hidden;
        transition: all 0.3s ease;
        
        &.canvas-mobile {
          max-width: 375px;
          box-shadow: 0 2px 20px rgba(0, 0, 0, 0.12);
        }

        .form-cover {
          position: relative;
          height: 200px;
          display: flex;
          align-items: center;
          justify-content: center;
          overflow: hidden;

          .cover-content {
            .cover-icon {
              color: rgba(255, 255, 255, 0.9);
              filter: drop-shadow(0 2px 8px rgba(0, 0, 0, 0.15));
            }
          }

          .edit-cover-btn {
            position: absolute;
            top: 16px;
            right: 16px;
            color: #fff;
            background: rgba(0, 0, 0, 0.2);

            &:hover {
              background: rgba(0, 0, 0, 0.3);
            }
          }
        }

        .form-header-section {
          padding: 40px 40px 32px;
          background: #fff;

          .form-title-wrapper {
            margin-bottom: 20px;

            .form-title-input {
              :deep(.t-input) {
                border: none !important;
                box-shadow: none !important;
                background: transparent;
                padding: 0;
              }
              
              :deep(.t-input__inner) {
                font-size: 36px;
                font-weight: 700;
                padding: 0;
                color: #1f2329;
                line-height: 1.3;
                letter-spacing: -0.5px;
                border: none;
                background: transparent;
                transition: color 0.2s ease;
                
                &::placeholder {
                  color: #c9cdd4;
                  font-weight: 400;
                }
                
                &:hover {
                  color: #1f2329;
                }
                
                &:focus {
                  color: #1f2329;
                  outline: none;
                }
              }
            }
          }

          .form-description-wrapper {
            .form-description-input {
              :deep(.t-textarea) {
                border: none !important;
                box-shadow: none !important;
                background: transparent;
                padding: 0;
              }
              
              :deep(.t-textarea__inner) {
                padding: 0;
                font-size: 16px;
                color: #646a73;
                line-height: 1.6;
                resize: none;
                border: none;
                background: transparent;
                min-height: auto;
                
                &::placeholder {
                  color: #c9cdd4;
                }
                
                &:hover {
                  color: #646a73;
                }
                
                &:focus {
                  color: #1f2329;
                  outline: none;
                }
              }
            }
          }
        }

        .form-fields {
          padding: 16px 32px 50px; // 底部留出50px
          min-height: 300px;
          position: relative;

          .drop-indicator {
            position: absolute;
            left: 32px;
            right: 32px;
            height: 3px;
            background: linear-gradient(90deg, #0052d9 0%, #4a90e2 100%);
            border-radius: 2px;
            pointer-events: none;
            z-index: 100;
            box-shadow: 0 0 8px rgba(0, 82, 217, 0.5);
            animation: pulse 1s ease-in-out infinite;
            
            &::before {
              content: '';
              position: absolute;
              left: -6px;
              top: -3px;
              width: 9px;
              height: 9px;
              background: #0052d9;
              border-radius: 50%;
              box-shadow: 0 0 6px rgba(0, 82, 217, 0.6);
            }
            
            &::after {
              content: '';
              position: absolute;
              right: -6px;
              top: -3px;
              width: 9px;
              height: 9px;
              background: #0052d9;
              border-radius: 50%;
              box-shadow: 0 0 6px rgba(0, 82, 217, 0.6);
            }
          }
          
          @keyframes pulse {
            0%, 100% {
              opacity: 1;
            }
            50% {
              opacity: 0.6;
            }
          }

          .form-field-item {
            display: flex;
            gap: 12px;
            padding: 16px;
            margin-bottom: 16px;
            border: 2px solid transparent;
            background: #fafafa;
            border-radius: 8px;
            cursor: pointer;
            transition: all 0.2s;
            position: relative;
            
            &.is-dragging {
              opacity: 0.4;
              transform: scale(0.98);
            }

            &:hover {
              border-color: #d9e8ff;
              background: #f5f9ff;

              .field-actions {
                opacity: 1;
              }
            }

            &.active {
              border-color: #0052d9;
              background: #f5f9ff;
              box-shadow: 0 0 0 2px rgba(0, 82, 217, 0.1);

              .field-actions {
                opacity: 1;
              }
            }

            &.drag-over {
              border-color: #0052d9;
              border-style: dashed;
              background: #e6f0ff;
            }

            .field-drag-handle {
              display: flex;
              align-items: center;
              justify-content: center;
              width: 24px;
              cursor: move;
              user-select: none;
              
              &:hover .drag-handle-icon span {
                background-color: #0052d9;
              }
              
              .drag-handle-icon {
                display: grid;
                grid-template-columns: repeat(2, 3px);
                gap: 3px;
                padding: 4px;
                
                span {
                  width: 3px;
                  height: 3px;
                  background-color: #c9cdd4;
                  border-radius: 50%;
                  transition: background-color 0.2s;
                }
              }
            }

            .field-content {
              flex: 1;
              min-width: 0;

              .field-header {
                display: flex;
                align-items: center;
                justify-content: space-between;
                margin-bottom: 12px;

                .field-label {
                  font-size: 14px;
                  font-weight: 500;
                  color: #333;

                  .required-mark {
                    color: #e34d59;
                    margin-right: 4px;
                  }
                }

                .field-actions {
                  display: flex;
                  gap: 4px;
                  opacity: 0;
                  transition: opacity 0.2s;
                }
              }

              .field-input {
                :deep(.t-input),
                :deep(.t-textarea),
                :deep(.t-select),
                :deep(.t-date-picker),
                :deep(.t-input-number) {
                  pointer-events: none;
                }
              }
            }
          }

          .empty-canvas {
            display: flex;
            flex-direction: column;
            align-items: center;
            justify-content: center;
            min-height: 400px;
            color: #999;

            .empty-icon {
              margin-bottom: 16px;
              color: #ddd;
            }

            .empty-title {
              margin: 8px 0;
              font-size: 16px;
              font-weight: 500;
              color: #666;
            }

            .empty-desc {
              font-size: 14px;
              color: #999;
            }
          }
        }
      }
    }

    // 右侧配置区
    .field-config-panel {
      width: 320px;
      background: #fff;
      border-left: 1px solid #e8e8e8;
      overflow-y: auto;
      flex-shrink: 0;

      .panel-header {
        display: flex;
        align-items: center;
        justify-content: space-between;
        padding: 16px;
        border-bottom: 1px solid #e8e8e8;
        font-weight: 500;
        position: sticky;
        top: 0;
        background: #fff;
        z-index: 10;
      }

      .panel-content {
        padding: 16px;

        .config-item {
          margin-bottom: 20px;

          .config-label {
            font-size: 13px;
            color: #666;
            margin-bottom: 8px;
          }

          .options-list {
            .option-item {
              display: flex;
              gap: 8px;
              margin-bottom: 8px;
            }
          }

          .range-inputs {
            display: flex;
            align-items: center;
            gap: 8px;
          }
        }
      }
    }
  }
}

// 统计弹窗样式
.stats-content {
  .stats-overview {
    display: grid;
    grid-template-columns: repeat(2, 1fr);
    gap: 16px;
    margin-bottom: 24px;

    .stat-card {
      display: flex;
      align-items: center;
      gap: 16px;
      padding: 20px;
      background: #fafafa;
      border-radius: 8px;

      .stat-icon {
        width: 48px;
        height: 48px;
        display: flex;
        align-items: center;
        justify-content: center;
        border-radius: 8px;
      }

      .stat-info {
        flex: 1;

        .stat-value {
          font-size: 24px;
          font-weight: 600;
          color: #333;
          margin-bottom: 4px;
        }

        .stat-label {
          font-size: 13px;
          color: #999;
        }
      }
    }
  }

  .stats-chart {
    margin-bottom: 24px;

    .chart-header {
      display: flex;
      align-items: center;
      justify-content: space-between;
      margin-bottom: 16px;

      .chart-title {
        font-size: 14px;
        font-weight: 500;
        color: #333;
      }
    }

    .chart-placeholder {
      display: flex;
      flex-direction: column;
      align-items: center;
      justify-content: center;
      height: 300px;
      background: #fafafa;
      border-radius: 8px;
      color: #999;

      p {
        margin-top: 12px;
        font-size: 14px;
      }
    }
  }

  .stats-actions {
    display: flex;
    justify-content: flex-end;
  }
}

// 主题弹窗样式
.theme-content {
  .theme-section {
    margin-bottom: 32px;

    &:last-child {
      margin-bottom: 0;
    }

    .section-title {
      margin-bottom: 16px;
      
      .title-text {
        display: block;
        font-size: 15px;
        font-weight: 600;
        color: #1f2329;
        margin-bottom: 4px;
      }
      
      .title-desc {
        display: block;
        font-size: 13px;
        color: #8a8f99;
      }
    }

    .gradient-grid {
      display: grid;
      grid-template-columns: repeat(3, 1fr);
      gap: 16px;

      .gradient-card {
        cursor: pointer;
        transition: all 0.2s;
        
        &:hover {
          transform: translateY(-2px);
        }
        
        &.selected {
          .gradient-preview {
            box-shadow: 0 0 0 3px var(--form-theme-color, #0052d9);
          }
        }

        .gradient-preview {
          height: 100px;
          border-radius: 8px;
          position: relative;
          transition: all 0.2s;
          box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
          
          .selected-badge {
            position: absolute;
            top: 8px;
            right: 8px;
            width: 28px;
            height: 28px;
            background: rgba(255, 255, 255, 0.95);
            border-radius: 50%;
            display: flex;
            align-items: center;
            justify-content: center;
            color: var(--form-theme-color, #0052d9);
            box-shadow: 0 2px 8px rgba(0, 0, 0, 0.15);
          }
        }
      }
    }

    .color-grid {
      display: grid;
      grid-template-columns: repeat(3, 1fr);
      gap: 16px;

      .color-card {
        cursor: pointer;
        text-align: center;
        transition: all 0.2s;
        
        &:hover {
          transform: translateY(-2px);
        }
        
        &.selected {
          .color-preview {
            box-shadow: 0 0 0 3px currentColor;
            transform: scale(1.05);
          }
        }

        .color-preview {
          width: 100%;
          height: 80px;
          border-radius: 8px;
          position: relative;
          margin-bottom: 8px;
          transition: all 0.2s;
          box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
          
          .selected-badge {
            position: absolute;
            top: 50%;
            left: 50%;
            transform: translate(-50%, -50%);
            width: 32px;
            height: 32px;
            background: rgba(255, 255, 255, 0.95);
            border-radius: 50%;
            display: flex;
            align-items: center;
            justify-content: center;
            color: currentColor;
            box-shadow: 0 2px 8px rgba(0, 0, 0, 0.2);
          }
        }
        
        .color-name {
          font-size: 13px;
          color: #646a73;
          font-weight: 500;
        }
      }
    }
  }
}

// 设置弹窗样式
.settings-content {
  .setting-item {
    display: flex;
    align-items: center;
    justify-content: space-between;
    padding: 16px 0;
    border-bottom: 1px solid #f0f0f0;

    &:last-child {
      border-bottom: none;
    }

    .setting-label {
      flex: 1;

      span:first-child {
        display: block;
        font-size: 14px;
        color: #333;
        margin-bottom: 4px;
      }

      .setting-desc {
        font-size: 12px;
        color: #999;
      }
    }
  }
}

// 分享弹窗样式
.share-content {
  .share-link-section {
    .link-display {
      margin-bottom: 24px;
    }

    .qrcode-section {
      margin-bottom: 24px;

      .qrcode-title {
        font-size: 14px;
        font-weight: 500;
        color: #333;
        margin-bottom: 12px;
      }

      .qrcode-placeholder {
        display: flex;
        flex-direction: column;
        align-items: center;
        justify-content: center;
        height: 200px;
        background: #fafafa;
        border-radius: 8px;
        border: 2px dashed #e8e8e8;
        color: #999;

        p {
          margin-top: 12px;
          font-size: 13px;
        }
      }
    }

    .share-options {
      display: flex;
      gap: 16px;
    }
  }

  .embed-section {
    :deep(.t-textarea) {
      margin-bottom: 16px;
    }
  }
}
</style>
