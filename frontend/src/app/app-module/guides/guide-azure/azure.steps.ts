import {Step} from '../shared/step';

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
    name: `Use the data collected in the previous steps to fill in the form as documented below.`,
    content: {
      id: 'stepContent5'
    }
  },
  {
    id: '6',
    name: `Click on the button shown below, to activate the UTMStack features related to this integration.`,
    content: {
      id: 'stepContent6'
    }
  }
];

