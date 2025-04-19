<script setup lang="ts">
import { ref, onMounted, h  } from 'vue'
import { Button } from '@/components/ui/button'
import {
    Dialog,
    DialogContent,
    DialogDescription,
    DialogFooter,
    DialogHeader,
    DialogTitle,
    DialogTrigger,
} from '@/components/ui/dialog'
import { checkConfig, sendConfig } from '@/lib/config'
import {
  Form,
  FormControl,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from '@/components/ui/form'
import { Input } from '@/components/ui/input'
import { toTypedSchema } from '@vee-validate/zod'
import * as z from 'zod'
import { toast } from 'vue-sonner'
import { Toaster } from '@/components/ui/sonner'


const formSchema = toTypedSchema(z.object({
  submit_flag_checker_time: z.number().min(0, 'Submit Flag Checker Time must be a positive number').default(120),
  host_flagchecker: z.string().min(1, 'Host Flag Checker is required').default(''),
  team_token: z.string().min(1, 'Team Token is required').default(''),
  max_flag_batch_size: z.number().min(1, 'Max Flag Batch Size must be a positive number').default(500),
  protocol: z.string().min(1, 'Protocol is required').default(''),
  base_url_server: z.string().url('Base URL Server must be a valid URL').default(''),
  submit_flag_server_time: z.number().min(0, 'Submit Flag Server Time must be a positive number').default(120),
  services: z.string().min(1, 'Services is required').default(''),
  range_ip_teams: z.string().min(1, 'IP Range for Teams is required').default(""),
  format_ip_teams: z.string().min(1, 'IP Format for Teams is required').default(''),
  my_team_ip: z.string().ip('My Team IP must be a valid IP address').default(''),
}));

const dialogOpen = ref(false)

function onSubmit(values: {string}) {
  const config = {
    configured:true,
    server: {
      team_token: values.team_token,
      host_flagchecker: values.host_flagchecker,
      protocol: values.protocol,
      max_flag_batch_size: values.max_flag_batch_size,
      submit_flag_checker_time: values.submit_flag_checker_time,
    },
    client: {
      base_url_server: values.base_url_server,
      submit_flag_server_time: values.submit_flag_server_time,
      services: [
        {name:"service1",port:80},
        {name:"service2",port:443}
      ],
      range_ip_teams: values.range_ip_teams,
      format_ip_teams: values.format_ip_teams,
      my_team_ip: values.my_team_ip,
    },
  };

  try{
    sendConfig(config);
    dialogOpen.value = false
    toast.success("Config sent successfully",{
      description: h('pre', { class: 'mt-2 w-[340px] rounded-md bg-slate-950 p-4' }, h('code', { class: 'text-white' }, JSON.stringify(config, null, 2))),
    })
  }catch(e){
    console.error(e)
    toast.error("Error sending config",{
      description: h('pre', { class: 'mt-2 w-[340px] rounded-md bg-slate-950 p-4' }, h('code', { class: 'text-white' }, JSON.stringify(e, null, 2))),
    })
  }
}



onMounted(async () => {
  if (!await checkConfig()) {
    dialogOpen.value = true
  }
})

</script>

<template>
<Form v-slot="{ handleSubmit }" class="space-y-4" as="" keep-values :validation-schema="formSchema">
    <Toaster rich-colors/>
  <Dialog :open="dialogOpen" @update:open="dialogOpen = $event">
    <DialogTrigger>
      <!-- Optional trigger button for dialog -->
      <!-- <Button variant="outline">Open Dialog</Button> -->
    </DialogTrigger>
    <DialogContent class="sm:max-w-3/6">
      <DialogHeader>
        <DialogTitle>Setup config</DialogTitle>
        <DialogDescription>
          Init the config by providing the necessary information.
        </DialogDescription>
      </DialogHeader>

      <form id="dialogForm" @submit="handleSubmit($event, onSubmit)">
        <!-- Server Configuration -->
        <h3 class="text-lg font-semibold mb-2 sm:col-span-2">Server Configuration</h3>
        <div class="grid grid-cols-1 sm:grid-cols-2 gap-6">
          <!-- Server Token -->
          <FormField v-slot="{ componentField }" name="team_token">
            <FormItem>
              <FormLabel>Team Token <span class="text-red-500">*</span></FormLabel>
              <FormControl>
                <Input type="text" placeholder="your-team-token" v-bind="componentField" />
              </FormControl>
              <FormMessage />
            </FormItem>
          </FormField>

          <!-- Host Flag Checker -->
          <FormField v-slot="{ componentField }" name="host_flagchecker">
            <FormItem>
              <FormLabel>Host Flag Checker <span class="text-red-500">*</span></FormLabel>
              <FormControl>
                <Input type="text" placeholder="flagchecker.example.com" v-bind="componentField" />
              </FormControl>
              <FormMessage  />
            </FormItem>
          </FormField>

          <!-- Protocol -->
          <FormField v-slot="{ componentField }" name="protocol">
            <FormItem>
              <FormLabel>Flag checker Protocol (.so) <span class="text-red-500">*</span></FormLabel>
              <FormControl>
                <Input type="text" placeholder="The name of shared object without extension" v-bind="componentField" />
              </FormControl>
              <FormMessage />
            </FormItem>
          </FormField>

          <!-- Submit Flag Checker Time -->
          <FormField v-slot="{ componentField }" name="submit_flag_checker_time">
            <FormItem>
              <FormLabel>Submit Flag Checker Time (seconds)</FormLabel>
              <FormControl>
                <Input type="number" placeholder="120" v-bind="componentField" />
              </FormControl>
              <FormMessage />
            </FormItem>
          </FormField>

          <!-- Max Flag Batch Size -->
          <FormField v-slot="{ componentField }" name="max_flag_batch_size">
            <FormItem>
              <FormLabel>Max Flag Batch Size</FormLabel>
              <FormControl>
                <Input type="number" placeholder="500" v-bind="componentField" />
              </FormControl>
              <FormMessage  />
            </FormItem>
          </FormField>

          <!-- Client Configuration -->
          <h3 class="text-lg font-semibold mt-4 mb-2 sm:col-span-2">Client Configuration</h3>

          <!-- Base URL Server -->
          <FormField v-slot="{ componentField }" name="base_url_server">
            <FormItem>
              <FormLabel>Base URL Server <span class="text-red-500">*</span></FormLabel>
              <FormControl>
                <Input type="url" placeholder="http://localhost:8080" v-bind="componentField" />
              </FormControl>
              <FormMessage  />
            </FormItem>
          </FormField>

          <!-- Services -->
          <FormField v-slot="{ componentField }" name="services">
            <FormItem>
              <FormLabel>Services (comma separated) <span class="text-red-500">*</span></FormLabel>
              <FormControl>
                <Input type="text" placeholder="service1,service2" v-bind="componentField" />
              </FormControl>
              <FormMessage  />
            </FormItem>
          </FormField>

          <!-- IP Range -->
          <FormField v-slot="{ componentField }" name="range_ip_teams">
            <FormItem>
              <FormLabel>IP Range for Teams <span class="text-red-500">*</span></FormLabel>
              <FormControl>
                <Input type="text" placeholder="10" v-bind="componentField" />
              </FormControl>
              <FormMessage  />
            </FormItem>
          </FormField>

          <!-- IP Format -->
          <FormField v-slot="{ componentField }" name="format_ip_teams">
            <FormItem>
              <FormLabel>IP Format for Teams <span class="text-red-500">*</span></FormLabel>
              <FormControl>
                <Input type="text" placeholder="10.0.0.{}" v-bind="componentField" />
              </FormControl>
              <FormMessage />
            </FormItem>
          </FormField>

          <!-- Team IP -->
          <FormField v-slot="{ componentField }" name="my_team_ip">
            <FormItem>
              <FormLabel>My Team IP <span class="text-red-500">*</span></FormLabel>
              <FormControl>
                <Input type="text" placeholder="192.168.1.1" v-bind="componentField" />
              </FormControl>
              <FormMessage />
            </FormItem>
          </FormField>

          <!-- Submit Flag Server Time -->
          <FormField v-slot="{ componentField }" name="submit_flag_server_time">
            <FormItem>
              <FormLabel>Submit Flag Server Time (seconds)</FormLabel>
              <FormControl>
                <Input type="number" placeholder="120" v-bind="componentField" />
              </FormControl>
              <FormMessage  />
            </FormItem>
          </FormField>
        </div>
        </form>

        <DialogFooter>
          <Button type="submit" form="dialogForm">
            Save changes
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  </Form>
</template>
