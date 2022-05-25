import { Id } from '../../definitions/key-types';

export class OrganizationSetting {
    public name!: string;
    public description!: string;
    public legal_notice!: string;
    public privacy_policy!: string;
    public login_text!: string;
    public theme_id!: Id; // (theme/theme_for_organization_id);
    public url!: string;
    public reset_password_verbose_errors!: boolean;
    public enable_electronic_voting!: boolean;
    public enable_chat!: boolean;
    public limit_of_meetings!: number;
    public limit_of_users!: number;
}

export class Organization {
    public static COLLECTION = `organization`;

    public name!: string;
    public description!: string;

    public committee_ids!: Id[]; // (committee/organization_id)[];
    public resource_ids!: Id[]; // (resource/organization_id)[];
    public organization_tag_ids!: Id[]; // (organization_tag/organization_id)[];
    public theme_ids!: Id[]; // (theme/organization_id);
    public active_meeting_ids!: Id[]; // (meeting/is_active_in_organization_id)[];
    public archived_meeting_ids!: Id[]; // (meeting/is_archived_in_organization_id)[];
    public template_meeting_ids!: Id[]; // (meeting/template_for_organization_id)[];
}

export interface Organization extends OrganizationSetting {}
