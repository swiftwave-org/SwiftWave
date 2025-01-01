<script setup>
import FilledButton from '@/views/components/FilledButton.vue';
import ModalDialog from '@/views/components/ModalDialog.vue';
import { reactive, ref } from 'vue';
import { useMutation } from '@vue/apollo-composable';
import gql from 'graphql-tag';
import { toast } from 'vue-sonner'
import { preventSpaceInput } from '@/vendor/utils.js';

const props = defineProps({
    callbackOnCreate: {
        type: Function,
        required: false,
        default: () => { }
    },
    callbackOnPop: {
        type: Function,
        required: false,
        default: () => { }
    }
});

const isModalOpen = ref(false);
const isInvalidDomainName = ref(false);
const domainRegex = /^((?:[-a-z0-9]+\.)+)([a-z]{2,63})\/?$/i;

const openModal = () => {
    newDomainDetails.name = '';
    isInvalidDomainName.value = false;
    isModalOpen.value = true;
};
const closeModal = () => {
    isModalOpen.value = false;
    props.callbackOnPop();
};

// Register Domain state
const newDomainDetails = reactive({
    name: ''
});

const {
    mutate: registerDomainRaw,
    loading: isDomainRegistering,
    onDone: onDomainRegisterSuccess,
    onError: onDomainRegisterFail
} = useMutation(
    gql`
    mutation ($input: DomainInput!) {
      addDomain(input: $input) {
        id
        name
      }
    }
  `,
    {
        variables: {
            input: newDomainDetails
        }
    }
);

onDomainRegisterSuccess((r) => {
    closeModal();
    newDomainDetails.name = '';
    isInvalidDomainName.value = false;
    toast.success('Domain registered successfully');
    props.callbackOnCreate(r.data.addDomain.id);
});

onDomainRegisterFail((err) => {
    toast.error(err.message);
});

const registerDomain = async () => {
    if (newDomainDetails.name === '') {
        isInvalidDomainName.value = true;
        return;
    }
    // validate name
    if (!domainRegex.test(newDomainDetails.name)) {
        isInvalidDomainName.value = true;
        return;
    }
    registerDomainRaw();
};

defineExpose({
    openModal,
    closeModal
});
</script>

<template>
    <teleport to="body">
        <ModalDialog :close-modal="closeModal" :is-open="isModalOpen">
            <template v-slot:header>Register New Domain</template>
            <template v-slot:body>
                Enter the domain or subdomain name you want to register.
                <form @submit.prevent="">
                    <!--  Name Field   -->
                    <div class="mt-4">
                        <label class="block text-sm font-medium text-gray-700" for="name">
                            Domain Name (example: example.com)
                        </label>
                        <div class="mt-1">
                            <input id="name" v-model="newDomainDetails.name" @keydown="preventSpaceInput"
                                autocomplete="off"
                                class="block w-full rounded-md border-gray-300 shadow-sm focus:border-primary-500 focus:ring-primary-500 sm:text-sm"
                                name="name" placeholder="example.com or test.example.com" type="text" />
                        </div>
                    </div>
                    <p v-if="isInvalidDomainName" class="mt-0.5 text-sm font-normal text-danger-500">
                        Provide a valid domain name
                    </p>
                </form>
            </template>
            <template v-slot:footer>
                <FilledButton :click="registerDomain" :loading="isDomainRegistering" type="primary">Register
                </FilledButton>
            </template>
        </ModalDialog>
    </teleport>
</template>

<style scoped></style>
