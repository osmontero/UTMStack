import { Step } from '../shared/step';

export const AZURE_STEPS: Step[] = [
  {
    id: '1',
    name: `Create a Log Analytics Workspace, using the official Azure documentation
      <a href="https://learn.microsoft.com/en-us/azure/azure-monitor/logs/quick-create-workspace" target="_blank">Create a Log Analytics workspace</a>.
      After creation, copy the <strong>Workspace ID</strong> from Overview. It will be needed in the next steps.`
  },
  {
    id: '2',
    name: `Register an Application in Microsoft Entra ID, using the official Azure documentation:
      <a href="https://learn.microsoft.com/en-us/entra/identity-platform/howto-app-registration-portal" target="_blank">How to register an Application in Microsoft Entra ID</a>.
      From the Overview, copy the <strong>Application (client) ID</strong> and the <strong>Directory (tenant) ID</strong>. You will need them later.`
  },
  {
    id: '3',
    name: `Generate a Client Secret for your Application using the official Azure documentation.
      <a href="https://learn.microsoft.com/en-us/entra/identity-platform/howto-add-app-credentials" target="_blank">Add and manage app credentials in Microsoft Entra ID</a>.
      Copy the <strong>Client Secret Value</strong>. You will need it to configure your tenant.`
  },
  {
    id: '4',
    name: `Grant workspace permissions. In the previously created workspace → Access Control (IAM) → Add → Add role assignment →
      select <strong>Log Analytics Reader</strong> as Role → select the <strong>Application</strong> created in Step 2 as Member.
      If you need more information, use the official Azure documentation:
      <a href="https://learn.microsoft.com/en-us/azure/azure-monitor/logs/manage-access" target="_blank">Manage access to Log Analytics Workspace</a>.`
  },
  {
    id: '5',
    name: `Enable sending logs from the subscription (Activity Logs). In <strong>Monitor</strong>, from the left menu, choose <strong>Activity log</strong>.
      Go to <strong>Export Activity Logs</strong>. Choose your subscription, and click on <strong>+Add diagnostic setting</strong>.
      Enter a name in <strong>Diagnostic setting name</strong>. In <strong>Logs</strong>, select the categories to send.
      In <strong>Destination details</strong>, choose <strong>Send to Log Analytics workspace</strong>. Select your <strong>Subscription</strong> and the <strong>Workspace</strong> from Step 1.
      Click <strong>Save</strong>. For more info:
      <a href="https://learn.microsoft.com/en-us/azure/azure-monitor/essentials/diagnostic-settings" target="_blank">Create diagnostic settings in Azure Monitor</a>.`
  },
  {
    id: '6',
    name: `Create a Data Collection Rule to route logs to the Workspace from Step 1.
      In <strong>Monitor</strong> → <strong>Settings</strong> → <strong>Data Collection Rules</strong> → <strong>Create</strong>.
      In <strong>Basics</strong>: enter Rule Name, Subscription, Resource Group, Region. Set Platform Type (recommended: All) → <strong>Next: Resources ></strong>.
      In <strong>Resources</strong>: click <strong>+Add resources</strong>, choose Scope, click <strong>Apply</strong> → <strong>Next: Collect and deliver ></strong>.
      In <strong>Collect and deliver</strong>: click <strong>+Add data source</strong>, choose source type → <strong>Next: Destination ></strong>.
      In <strong>Destination</strong>: click <strong>+Add destination</strong>, choose <strong>Azure Monitor Logs</strong>, your Subscription, and the Workspace from Step 1.
      Click <strong>Add data source</strong>. Repeat to add more sources. Click <strong>Review + create</strong>.
      More info:
      <a href="https://learn.microsoft.com/en-us/azure/azure-monitor/agents/data-collection-rule-overview" target="_blank">Data collection rules in Azure Monitor</a>.`
  },
  {
    id: '7',
    name: `Use the data collected in the previous steps to fill in the form as documented below.`,
    content: {
      id: 'stepContent7'
    }
  },
  {
    id: '8',
    name: `Click on the button shown below, to activate the UTMStack features related to this integration.`,
    content: {
      id: 'stepContent8'
    }
  }
];
