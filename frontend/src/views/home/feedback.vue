<template>
  <div class="feedback-page">
    <!-- 学习教程区域 -->
    <div class="tutorial-section">
      <div class="section-title">学习教程</div>
      <div class="tutorial-list">
        <div
          v-for="tutorial in tutorialList"
          :key="tutorial.id"
          class="tutorial-card"
          @click="handlePlayVideo(tutorial)"
        >
          <div class="tutorial-cover">
            <img :src="tutorial.cover" :alt="tutorial.title" />
            <div class="play-icon">
              <t-icon name="play-circle" size="48px" />
            </div>
          </div>
          <div class="tutorial-info">
            <div class="tutorial-title">{{ tutorial.title }}</div>
            <div class="tutorial-desc">{{ tutorial.desc }}</div>
          </div>
        </div>
      </div>
    </div>

    <!-- 问题与反馈区域 -->
    <div class="qa-section">
      <div class="section-title">问题与反馈</div>
      <div class="qa-container">
        <t-collapse :default-value="['1']" expand-icon-placement="right">
          <t-collapse-panel
            v-for="qa in qaList"
            :key="qa.id"
            :value="qa.id"
            :header="qa.question"
          >
            <div class="qa-answer" v-html="qa.answer"></div>
          </t-collapse-panel>
        </t-collapse>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { MessagePlugin } from 'tdesign-vue-next'

// 教程数据接口
interface Tutorial {
  id: string
  title: string
  desc: string
  cover: string
  videoUrl?: string
}

// QA数据接口
interface QAItem {
  id: string
  question: string
  answer: string
  category?: string
}

// 视频教程列表
const tutorialList = ref<Tutorial[]>([
  {
    id: '1',
    title: '什么是多维表？',
    desc: '基础功能介绍😊',
    cover: 'https://images.unsplash.com/photo-1460925895917-afdab827c52f?w=400&h=250&fit=crop',
  },
  {
    id: '2',
    title: '什么是视图？',
    desc: '基础高频问答展示✔️',
    cover: 'https://images.unsplash.com/photo-1551288049-bebda4e38f71?w=400&h=250&fit=crop',
  },
  {
    id: '3',
    title: '用多维表高效协作',
    desc: '帮非技术背景员工💡',
    cover: 'https://images.unsplash.com/photo-1522071820081-009f0129c71c?w=400&h=250&fit=crop',
  },
  {
    id: '4',
    title: '模板使用教学',
    desc: '快速掌握模板使用😎',
    cover: 'https://images.unsplash.com/photo-1516321318423-f06f85e504b3?w=400&h=250&fit=crop',
  },
])

// 常见问题列表
const qaList = ref<QAItem[]>([
  {
    id: '1',
    category: '基础功能',
    question: '如何创建一个新的多维表？',
    answer: `
      <p>创建新的多维表非常简单，请按照以下步骤操作：</p>
      <ol>
        <li>点击首页右上角的「<strong>+ 新建</strong>」按钮</li>
        <li>选择「从空白创建」或「从模板创建」</li>
        <li>如果选择从空白创建，输入表格名称并点击确认</li>
        <li>如果选择从模板创建，选择合适的模板后点击「使用」按钮</li>
        <li>系统会自动跳转到新创建的表格编辑页面</li>
      </ol>
      <p><strong>提示：</strong>建议首次使用时从模板创建，可以更快地了解多维表的功能。</p>
    `,
  },
  {
    id: '2',
    category: '数据管理',
    question: '如何导入已有的数据到多维表？',
    answer: `
      <p>支持多种数据导入方式：</p>
      <ul>
        <li><strong>Excel/CSV 导入：</strong>点击表格右上角「导入」按钮，选择 Excel 或 CSV 文件上传</li>
        <li><strong>复制粘贴：</strong>从 Excel 或其他表格软件复制数据，直接粘贴到多维表中</li>
        <li><strong>API 导入：</strong>通过 API 接口批量导入数据（需要专业版或以上权限）</li>
      </ul>
      <p><strong>注意事项：</strong></p>
      <ul>
        <li>单次导入数据量不超过 10000 行（基础版）</li>
        <li>导入时会自动识别数据类型，也可以手动调整</li>
        <li>建议先预览数据，确认无误后再正式导入</li>
      </ul>
    `,
  },
  {
    id: '3',
    category: '视图管理',
    question: '什么是视图？如何创建和使用视图？',
    answer: `
      <p><strong>视图</strong>是对同一份数据的不同展示方式，可以根据不同的使用场景创建多个视图。</p>
      <p><strong>常见视图类型：</strong></p>
      <ul>
        <li><strong>表格视图：</strong>传统的行列表格形式，适合数据录入和编辑</li>
        <li><strong>看板视图：</strong>卡片式展示，适合任务管理和项目追踪</li>
        <li><strong>日历视图：</strong>时间轴展示，适合日程安排和事件管理</li>
        <li><strong>画廊视图：</strong>图片墙展示，适合产品目录和作品展示</li>
      </ul>
      <p><strong>创建视图步骤：</strong></p>
      <ol>
        <li>点击视图切换区域的「+ 添加视图」按钮</li>
        <li>选择视图类型</li>
        <li>设置视图名称和显示字段</li>
        <li>配置筛选、排序和分组条件</li>
        <li>保存视图即可使用</li>
      </ol>
    `,
  },
  {
    id: '4',
    category: '协作功能',
    question: '如何邀请团队成员协作？',
    answer: `
      <p>多维表支持强大的团队协作功能：</p>
      <ol>
        <li>点击表格右上角的「分享」按钮</li>
        <li>选择「邀请成员」选项</li>
        <li>输入成员的邮箱地址（支持批量邀请，用逗号分隔）</li>
        <li>设置成员权限：
          <ul>
            <li><strong>可编辑：</strong>可以查看和修改数据</li>
            <li><strong>仅查看：</strong>只能查看数据，不能修改</li>
            <li><strong>管理员：</strong>拥有所有权限，包括成员管理</li>
          </ul>
        </li>
        <li>点击「发送邀请」，系统会发送邮件通知</li>
      </ol>
      <p><strong>权限管理：</strong>可以随时在「成员管理」中调整成员权限或移除成员。</p>
    `,
  },
  {
    id: '5',
    category: '高级功能',
    question: '如何使用公式和函数？',
    answer: `
      <p>多维表支持丰富的公式和函数功能，类似于 Excel：</p>
      <p><strong>常用函数分类：</strong></p>
      <ul>
        <li><strong>数学函数：</strong>SUM（求和）、AVERAGE（平均值）、MAX（最大值）、MIN（最小值）</li>
        <li><strong>文本函数：</strong>CONCATENATE（连接文本）、LEFT（提取左侧字符）、RIGHT（提取右侧字符）</li>
        <li><strong>日期函数：</strong>NOW（当前时间）、DATEADD（日期加减）、DATEDIF（日期差）</li>
        <li><strong>逻辑函数：</strong>IF（条件判断）、AND（与）、OR（或）、NOT（非）</li>
        <li><strong>查找函数：</strong>LOOKUP（查找引用）、ROLLUP（汇总统计）</li>
      </ul>
      <p><strong>使用方法：</strong></p>
      <ol>
        <li>在需要使用公式的字段中，点击字段类型选择「公式」</li>
        <li>在公式编辑器中输入公式，以 = 开头</li>
        <li>可以引用其他字段，使用 {字段名} 的格式</li>
        <li>示例：<code>=IF({状态}="已完成", "✅", "⏳")</code></li>
      </ol>
    `,
  },
  {
    id: '6',
    category: '数据安全',
    question: '我的数据安全吗？会丢失吗？',
    answer: `
      <p>我们非常重视数据安全，采取了多重保障措施：</p>
      <ul>
        <li><strong>数据加密：</strong>传输过程使用 HTTPS 加密，存储使用 AES-256 加密</li>
        <li><strong>自动备份：</strong>每小时自动备份一次，保留最近 30 天的备份数据</li>
        <li><strong>版本控制：</strong>记录所有数据变更历史，可随时恢复到历史版本</li>
        <li><strong>权限控制：</strong>细粒度的权限管理，确保数据只有授权人员可以访问</li>
        <li><strong>异地容灾：</strong>数据存储在多个地理位置的服务器上，确保高可用性</li>
      </ul>
      <p><strong>数据恢复：</strong>如果误删除数据，可以在「回收站」中找回，或联系客服协助恢复。</p>
    `,
  },
  {
    id: '7',
    category: '性能优化',
    question: '表格数据量很大时会卡顿吗？如何优化？',
    answer: `
      <p>系统采用了多种优化技术，确保大数据量下的流畅体验：</p>
      <ul>
        <li><strong>虚拟滚动：</strong>只渲染可见区域的数据，支持百万级数据流畅滚动</li>
        <li><strong>按需加载：</strong>数据分页加载，避免一次性加载过多数据</li>
        <li><strong>智能索引：</strong>自动建立索引，加快查询和筛选速度</li>
        <li><strong>缓存机制：</strong>频繁访问的数据会被缓存，减少服务器请求</li>
      </ul>
      <p><strong>优化建议：</strong></p>
      <ol>
        <li>合理使用筛选和视图，只显示需要的数据</li>
        <li>避免在单个表格中存储过多字段（建议不超过 50 个）</li>
        <li>定期归档历史数据，保持表格整洁</li>
        <li>使用链接字段代替重复数据，提高数据复用性</li>
      </ol>
    `,
  },
  {
    id: '8',
    category: '账号与权限',
    question: '如何升级到专业版或旗舰版？',
    answer: `
      <p>升级到付费版本可以解锁更多高级功能和更大的使用限额。</p>
      <p><strong>升级步骤：</strong></p>
      <ol>
        <li>进入「个人设置」页面</li>
        <li>在「用户权益」区域点击「升级」按钮</li>
        <li>选择适合的套餐（专业版 199元/年 或 旗舰版 999元/年）</li>
        <li>选择支付方式（支持微信、支付宝、企业转账）</li>
        <li>完成支付后即刻生效</li>
      </ol>
      <p><strong>版本对比：</strong></p>
      <table style="width: 100%; border-collapse: collapse; margin-top: 10px;">
        <tr>
          <th style="border: 1px solid #ddd; padding: 8px;">功能</th>
          <th style="border: 1px solid #ddd; padding: 8px;">基础版</th>
          <th style="border: 1px solid #ddd; padding: 8px;">专业版</th>
          <th style="border: 1px solid #ddd; padding: 8px;">旗舰版</th>
        </tr>
        <tr>
          <td style="border: 1px solid #ddd; padding: 8px;">表格数量</td>
          <td style="border: 1px solid #ddd; padding: 8px;">100个</td>
          <td style="border: 1px solid #ddd; padding: 8px;">300个</td>
          <td style="border: 1px solid #ddd; padding: 8px;">1000个</td>
        </tr>
        <tr>
          <td style="border: 1px solid #ddd; padding: 8px;">单表行数</td>
          <td style="border: 1px solid #ddd; padding: 8px;">2000行</td>
          <td style="border: 1px solid #ddd; padding: 8px;">10000行</td>
          <td style="border: 1px solid #ddd; padding: 8px;">20000行</td>
        </tr>
        <tr>
          <td style="border: 1px solid #ddd; padding: 8px;">团队人数</td>
          <td style="border: 1px solid #ddd; padding: 8px;">10人</td>
          <td style="border: 1px solid #ddd; padding: 8px;">100人</td>
          <td style="border: 1px solid #ddd; padding: 8px;">1000人</td>
        </tr>
      </table>
    `,
  },
])

// 播放视频
const handlePlayVideo = (tutorial: Tutorial) => {
  MessagePlugin.info(`播放视频: ${tutorial.title}`)
  // TODO: 打开视频播放弹窗或跳转到视频页面
  console.log('播放视频:', tutorial)
}
</script>

<style scoped lang="less">
.feedback-page {
  background: #ffffff;
  border-radius: 8px;
  padding: 27px;
  min-height: 100%;

  .section-title {
    font-size: 18px;
    font-weight: 500;
    color: #000000;
    margin-bottom: 20px;
  }

  // 学习教程区域
  .tutorial-section {
    margin-bottom: 40px;

    .tutorial-list {
      display: grid;
      grid-template-columns: repeat(auto-fill, minmax(240px, 1fr));
      gap: 16px;

      .tutorial-card {
        border: 1px solid #e8e8e8;
        border-radius: 8px;
        overflow: hidden;
        cursor: pointer;
        transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);

        &:hover {
          transform: translateY(-4px);
          box-shadow: 0 8px 24px rgba(74, 127, 247, 0.2);
          border-color: #4a7ff7;

          .play-icon {
            opacity: 1;
            transform: scale(1.1);
          }
        }

        .tutorial-cover {
          position: relative;
          width: 100%;
          height: 150px;
          overflow: hidden;
          background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);

          img {
            width: 100%;
            height: 100%;
            object-fit: cover;
          }

          .play-icon {
            position: absolute;
            top: 50%;
            left: 50%;
            transform: translate(-50%, -50%);
            color: #ffffff;
            opacity: 0.8;
            transition: all 0.3s;
          }
        }

        .tutorial-info {
          padding: 12px 14px;

          .tutorial-title {
            font-size: 15px;
            font-weight: 500;
            color: #000000;
            margin-bottom: 6px;
            overflow: hidden;
            text-overflow: ellipsis;
            white-space: nowrap;
          }

          .tutorial-desc {
            font-size: 13px;
            color: #666;
            overflow: hidden;
            text-overflow: ellipsis;
            white-space: nowrap;
          }
        }
      }
    }
  }

  // QA区域
  .qa-section {
    .qa-container {
      max-width: 900px;

      :deep(.t-collapse) {
        border: none;
      }

      :deep(.t-collapse-panel) {
        margin-bottom: 12px;
        border: 1px solid #e8e8e8;
        border-radius: 8px;
        overflow: hidden;

        &:last-child {
          margin-bottom: 0;
        }
      }

      :deep(.t-collapse-panel__header) {
        padding: 16px 20px;
        font-size: 15px;
        font-weight: 500;
        color: #000000;
        background: #fafafa;
        transition: all 0.3s;

        &:hover {
          background: #f0f0f0;
          color: #4a7ff7;
        }
      }

      :deep(.t-collapse-panel__body) {
        padding: 20px;
        background: #ffffff;
      }

      :deep(.t-collapse-panel--active) {
        border-color: #4a7ff7;

        .t-collapse-panel__header {
          background: #e8f0ff;
          color: #4a7ff7;
        }
      }

      .qa-answer {
        font-size: 14px;
        line-height: 1.8;
        color: #333;

        :deep(p) {
          margin: 0 0 12px 0;

          &:last-child {
            margin-bottom: 0;
          }
        }

        :deep(ul), :deep(ol) {
          margin: 12px 0;
          padding-left: 24px;
        }

        :deep(li) {
          margin: 8px 0;
          line-height: 1.6;
        }

        :deep(strong) {
          color: #000;
          font-weight: 600;
        }

        :deep(code) {
          padding: 2px 6px;
          background: #f5f5f5;
          border-radius: 3px;
          font-family: 'Monaco', 'Courier New', monospace;
          font-size: 13px;
          color: #e83e8c;
        }

        :deep(table) {
          border-collapse: collapse;
          margin: 12px 0;
          font-size: 13px;

          th, td {
            border: 1px solid #ddd;
            padding: 8px;
            text-align: left;
          }

          th {
            background: #f5f5f5;
            font-weight: 600;
          }

          tr:hover {
            background: #fafafa;
          }
        }
      }
    }
  }
}
</style>
